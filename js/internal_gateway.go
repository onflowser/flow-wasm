package js

import (
	"context"
	"encoding/json"
	"github.com/onflow/cadence"
	sdk "github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flowkit/v2/gateway"
	"syscall/js"
)

// Mapping as defined in https://github.com/onflow/fcl-js/blob/9c7873140015c9d1e28712aed93c56654f656639/packages/transport-grpc/src/send-get-account.js#L16-L28
var hashAlgoToJsIndex = map[crypto.HashAlgorithm]int{
	crypto.SHA2_256: 1,
	crypto.SHA2_384: 2,
	crypto.SHA3_256: 3,
	crypto.SHA3_384: 4,
	crypto.KMAC128:  5,
}

// Mapping as defined in https://github.com/onflow/fcl-js/blob/9c7873140015c9d1e28712aed93c56654f656639/packages/transport-grpc/src/send-get-account.js#L16-L28
var signAlgoToJsIndex = map[crypto.SignatureAlgorithm]int{
	crypto.ECDSA_P256:      1,
	crypto.ECDSA_secp256k1: 2,
	crypto.BLS_BLS12_381:   3,
}

type InternalGateway struct {
	emulator *gateway.EmulatorGateway
	target   js.Value
}

func NewInternalGateway(emulator *gateway.EmulatorGateway) *InternalGateway {
	target := js.Global().Get("Object").New()

	gtw := &InternalGateway{
		emulator,
		target,
	}

	target.Set("getAccount", js.FuncOf(gtw.getAccount))
	target.Set("getLatestBlock", js.FuncOf(gtw.getLatestBlock))
	target.Set("getBlockById", js.FuncOf(gtw.getBlockByID))
	target.Set("getBlockByHeight", js.FuncOf(gtw.getBlockByHeight))
	target.Set("getTransactionsByBlockId", js.FuncOf(gtw.getTransactionsByBlockID))
	target.Set("getTransaction", js.FuncOf(gtw.getTransaction))
	target.Set("getCollection", js.FuncOf(gtw.getCollection))
	target.Set("sendSignedTransaction", js.FuncOf(gtw.sendSignedTransaction))
	target.Set("getNetworkParameters", js.FuncOf(gtw.getNetworkParameters))

	return gtw
}

func (g *InternalGateway) JsValue() js.Value {
	return g.target
}

func (g *InternalGateway) getAccount(this js.Value, args []js.Value) interface{} {
	account, err := g.emulator.GetAccount(context.Background(), sdk.HexToAddress(args[0].String()))

	if err != nil {
		panic(err)
	}

	return serializeAccount(account)
}

func serializeAccount(account *sdk.Account) interface{} {
	serializedContracts := make(map[string]interface{})
	for key, value := range account.Contracts {
		serializedContracts[key] = string(value)
	}

	serializedKeys := make([]interface{}, 0)
	for _, value := range account.Keys {
		serializedKeys = append(serializedKeys, map[string]interface{}{
			"index":          value.Index,
			"publicKey":      value.PublicKey.String(),
			"signAlgo":       signAlgoToJsIndex[value.SigAlgo],
			"signAlgoString": value.SigAlgo.String(),
			"hashAlgo":       hashAlgoToJsIndex[value.HashAlgo],
			"hashAlgoString": value.HashAlgo.String(),
			"weight":         value.Weight,
			"sequenceNumber": value.SequenceNumber,
			"revoked":        value.Revoked,
		})
	}

	// https://developers.flow.com/tools/clients/fcl-js/api#accountobject
	return map[string]interface{}{
		"address":   account.Address.String(),
		"balance":   account.Balance,
		"contracts": serializedContracts,
		"keys":      serializedKeys,
		"code":      string(account.Code),
	}
}

type SendSignedTransactionRequest struct {
	Script           string `json:"script"`
	ReferenceBlockID string `json:"referenceBlockId"`
	GasLimit         uint64 `json:"gasLimit"`
}

func (g *InternalGateway) sendSignedTransaction(this js.Value, args []js.Value) interface{} {
	var request SendSignedTransactionRequest
	err := json.Unmarshal([]byte(args[0].String()), &request)

	inputTx := &sdk.Transaction{
		Script:             []byte(request.Script),
		Arguments:          [][]byte{},
		ReferenceBlockID:   sdk.HexToID(request.ReferenceBlockID),
		GasLimit:           request.GasLimit,
		ProposalKey:        sdk.ProposalKey{},
		Payer:              sdk.ServiceAddress(sdk.Emulator),
		Authorizers:        nil,
		PayloadSignatures:  nil,
		EnvelopeSignatures: nil,
	}

	outputTx, err := g.emulator.SendSignedTransaction(context.Background(), inputTx)

	if err != nil {
		panic(err)
	}

	return outputTx.ID().Hex()
}

func (g *InternalGateway) getTransaction(this js.Value, args []js.Value) interface{} {
	tx, err := g.emulator.GetTransaction(context.Background(), sdk.HexToID(args[0].String()))

	if err != nil {
		panic(err)
	}

	return serializeTransaction(tx)
}

func (g *InternalGateway) getTransactionsByBlockID(this js.Value, args []js.Value) interface{} {
	txs, err := g.emulator.GetTransactionsByBlockID(context.Background(), sdk.HexToID(args[0].String()))

	if err != nil {
		panic(err)
	}

	serializedTransactions := make([]interface{}, 0)
	for _, tx := range txs {
		serializedTransactions = append(serializedTransactions, serializeTransaction(tx))
	}

	return serializedTransactions
}

func serializeTransaction(tx *sdk.Transaction) map[string]interface{} {
	// https://developers.flow.com/tools/clients/fcl-js/api#proposalkeyobject
	serializedProposalKey := map[string]interface{}{
		"address":        tx.ProposalKey.Address,
		"keyIndex":       tx.ProposalKey.KeyIndex,
		"sequenceNumber": tx.ProposalKey.SequenceNumber,
	}

	serializedAuthorizers := make([]string, 0)
	for _, value := range tx.Authorizers {
		serializedAuthorizers = append(serializedAuthorizers, value.Hex())
	}

	serializedEnvelopeSignatures := make([]interface{}, 0)
	for _, value := range tx.EnvelopeSignatures {
		// https://developers.flow.com/tools/clients/fcl-js/api#signableobject
		serializedEnvelopeSignatures = append(serializedEnvelopeSignatures, map[string]interface{}{
			"addr":      value.Address.Hex(),
			"keyId":     value.KeyIndex,
			"signature": string(value.Signature),
		})
	}

	// https://developers.flow.com/tools/clients/fcl-js/api#transactionobject
	return map[string]interface{}{
		"id":                 tx.ID().Hex(),
		"authorizers":        serializedAuthorizers,
		"envelopeSignatures": serializedEnvelopeSignatures,
		"gasLimit":           tx.GasLimit,
		"payer":              tx.Payer.Hex(),
		"proposalKey":        serializedProposalKey,
		"referenceBlockId":   tx.ReferenceBlockID.Hex(),
		"script":             string(tx.Script),
	}
}

func (g *InternalGateway) getTransactionResultsByBlockID(ctx context.Context, blockID sdk.Identifier) ([]*sdk.TransactionResult, error) {
	//TODO implement me
	panic("implement me")
}

func (g *InternalGateway) getTransactionResult(ctx context.Context, identifier sdk.Identifier, b bool) (*sdk.TransactionResult, error) {
	//TODO implement me
	panic("implement me")
}

func (g *InternalGateway) executeScript(ctx context.Context, bytes []byte, values []cadence.Value) (cadence.Value, error) {
	//TODO implement me
	panic("implement me")
}

func (g *InternalGateway) executeScriptAtHeight(ctx context.Context, bytes []byte, values []cadence.Value, u uint64) (cadence.Value, error) {
	//TODO implement me
	panic("implement me")
}

func (g *InternalGateway) executeScriptAtID(ctx context.Context, bytes []byte, values []cadence.Value, identifier sdk.Identifier) (cadence.Value, error) {
	//TODO implement me
	panic("implement me")
}

func (g *InternalGateway) getLatestBlock(this js.Value, args []js.Value) interface{} {
	block, err := g.emulator.GetLatestBlock(context.Background())

	if err != nil {
		panic(err)
	}

	return serializeBlock(block)
}

func (g *InternalGateway) getBlockByHeight(this js.Value, args []js.Value) interface{} {
	block, err := g.emulator.GetBlockByHeight(context.Background(), uint64(args[0].Int()))

	if err != nil {
		panic(err)
	}

	return serializeBlock(block)
}

func (g *InternalGateway) getBlockByID(this js.Value, args []js.Value) interface{} {
	block, err := g.emulator.GetBlockByID(context.Background(), sdk.HexToID(args[0].String()))

	if err != nil {
		panic(err)
	}

	return serializeBlock(block)
}

func serializeBlock(block *sdk.Block) interface{} {
	serializedCollectionGuarantees := make([]interface{}, 0)

	for _, value := range block.CollectionGuarantees {
		serializedCollectionGuarantees = append(serializedCollectionGuarantees, map[string]interface{}{
			"collectionId": value.CollectionID.String(),
		})
	}

	serializedBlockSeals := make([]interface{}, 0)

	for _, value := range block.Seals {
		serializedBlockSeals = append(serializedBlockSeals, map[string]interface{}{
			"blockId":            value.BlockID.String(),
			"executionReceiptId": value.ExecutionReceiptID.String(),
		})
	}

	// https://developers.flow.com/tools/clients/fcl-js/api#blockobject
	return map[string]interface{}{
		"id":                   block.ID.String(),
		"parentId":             block.ParentID.String(),
		"height":               block.Height,
		"timestamp":            block.Timestamp.String(),
		"collectionGuarantees": serializedCollectionGuarantees,
		"blockSeals":           serializedBlockSeals,
		"signatures":           []interface{}{}, // Not implemented
	}
}

func (g *InternalGateway) getEvents(ctx context.Context, s string, u uint64, u2 uint64) ([]sdk.BlockEvents, error) {
	//TODO implement me
	panic("implement me")
}

func (g *InternalGateway) getCollection(this js.Value, args []js.Value) interface{} {
	collection, err := g.emulator.GetCollection(context.Background(), sdk.HexToID(args[0].String()))

	if err != nil {
		panic(err)
	}

	return serializeCollection(collection)
}

func serializeCollection(collection *sdk.Collection) interface{} {
	serializedTransactionIds := make([]string, 0)
	for _, value := range collection.TransactionIDs {
		serializedTransactionIds = append(serializedTransactionIds, value.Hex())
	}

	// https://developers.flow.com/tools/clients/fcl-js/api#collectionobject
	return map[string]interface{}{
		"id":             collection.ID().Hex(),
		"transactionIds": serializedTransactionIds,
	}
}

func (g *InternalGateway) getNetworkParameters(this js.Value, args []js.Value) interface{} {
	return map[string]interface{}{
		"chainId": sdk.Emulator.String(),
	}
}

func (g *InternalGateway) getLatestProtocolStateSnapshot(ctx context.Context) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (g *InternalGateway) ping() error {
	//TODO implement me
	panic("implement me")
}

func (g *InternalGateway) waitServer(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (g *InternalGateway) secureConnection() bool {
	//TODO implement me
	panic("implement me")
}
