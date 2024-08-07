package js

import (
	"context"
	"encoding/json"
	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	sdk "github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flowkit/v2/arguments"
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
	target.Set("getTransactionResultsByBlockId", js.FuncOf(gtw.getTransactionResultsByBlockID))
	target.Set("getTransactionResult", js.FuncOf(gtw.getTransactionResult))
	target.Set("executeScript", js.FuncOf(gtw.executeScript))
	target.Set("executeScriptAtHeight", js.FuncOf(gtw.executeScriptAtHeight))
	target.Set("executeScriptAtId", js.FuncOf(gtw.executeScriptAtID))

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
	Script           string   `json:"script"`
	ReferenceBlockID string   `json:"referenceBlockId"`
	GasLimit         uint64   `json:"gasLimit"`
	ArgsJSON         []string `json:"arguments"`
}

func (g *InternalGateway) sendSignedTransaction(this js.Value, args []js.Value) interface{} {
	var request SendSignedTransactionRequest
	err := json.Unmarshal([]byte(args[0].String()), &request)

	deserializedArgs := make([][]byte, 0)
	for _, arg := range request.ArgsJSON {
		deserializedArgs = append(deserializedArgs, []byte(arg))
	}

	inputTx := &sdk.Transaction{
		Script:             []byte(request.Script),
		Arguments:          deserializedArgs,
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
		"address":        tx.ProposalKey.Address.Hex(),
		"keyId":          tx.ProposalKey.KeyIndex,
		"sequenceNumber": tx.ProposalKey.SequenceNumber,
	}

	serializedAuthorizers := make([]interface{}, 0)
	for _, value := range tx.Authorizers {
		serializedAuthorizers = append(serializedAuthorizers, value.Hex())
	}

	serializedEnvelopeSignatures := make([]interface{}, 0)
	for _, value := range tx.EnvelopeSignatures {
		serializedEnvelopeSignatures = append(serializedEnvelopeSignatures, serializeSignature(value))
	}

	serializedPayloadSignatures := make([]interface{}, 0)
	for _, value := range tx.PayloadSignatures {
		serializedPayloadSignatures = append(serializedPayloadSignatures, serializeSignature(value))
	}

	// https://developers.flow.com/tools/clients/fcl-js/api#transactionobject
	return map[string]interface{}{
		"authorizers":        serializedAuthorizers,
		"envelopeSignatures": serializedEnvelopeSignatures,
		"payloadSignatures":  serializedPayloadSignatures,
		"gasLimit":           tx.GasLimit,
		"payer":              tx.Payer.Hex(),
		"proposalKey":        serializedProposalKey,
		"referenceBlockId":   tx.ReferenceBlockID.Hex(),
		"script":             string(tx.Script),
		"args":               []interface{}{}, // TODO: Implement
	}
}

func serializeSignature(sig sdk.TransactionSignature) interface{} {
	// https://developers.flow.com/tools/clients/fcl-js/api#signableobject
	return map[string]interface{}{
		"addr":      sig.Address.Hex(),
		"keyId":     sig.KeyIndex,
		"signature": string(sig.Signature),
	}
}

func (g *InternalGateway) getTransactionResultsByBlockID(this js.Value, args []js.Value) interface{} {
	results, err := g.emulator.GetTransactionResultsByBlockID(context.Background(), sdk.HexToID(args[0].String()))

	if err != nil {
		panic(err)
	}

	serializedResults := make([]interface{}, 0)
	for _, result := range results {
		serializedResults = append(serializedResults, serializeTransactionResult(result))
	}

	return serializedResults
}

func (g *InternalGateway) getTransactionResult(this js.Value, args []js.Value) interface{} {
	result, err := g.emulator.GetTransactionResult(context.Background(), sdk.HexToID(args[0].String()), false)

	if err != nil {
		panic(err)
	}

	return serializeTransactionResult(result)
}

func serializeTransactionResult(result *sdk.TransactionResult) interface{} {
	serializedEvents := make([]interface{}, 0)
	for _, event := range result.Events {
		// https://developers.flow.com/tools/clients/fcl-js/api#event-object
		serializedEvents = append(serializedEvents, map[string]interface{}{
			"type":             event.Type,
			"data":             event.Value.String(),
			"blockId":          result.BlockID.Hex(),
			"blockHeight":      result.BlockHeight,
			"blockTimestamp":   "", // TODO: Implement
			"transactionId":    event.TransactionID.Hex(),
			"transactionIndex": event.TransactionIndex,
			"eventIndex":       event.EventIndex,
		})
	}

	statusCode := 0
	if result.Error != nil {
		statusCode = 1
	}

	errorMessage := ""
	if result.Error != nil {
		errorMessage = result.Error.Error()
	}

	// https://developers.flow.com/tools/clients/fcl-js/api#transactionstatusobject
	return map[string]interface{}{
		"blockId":      result.BlockID.Hex(),
		"events":       serializedEvents,
		"status":       int(result.Status),
		"statusString": result.Status.String(),
		"errorMessage": errorMessage,
		"statusCode":   statusCode,
	}
}

type ExecuteScriptRequest struct {
	Script   string `json:"script"`
	ArgsJSON string `json:"arguments"`
}

func (g *InternalGateway) executeScript(this js.Value, args []js.Value) interface{} {
	var request ExecuteScriptRequest
	err := json.Unmarshal([]byte(args[0].String()), &request)

	if err != nil {
		panic(err)
	}

	var cadenceArgs []cadence.Value

	if request.ArgsJSON != "" {
		cadenceArgs, err = arguments.ParseJSON(request.ArgsJSON)

		if err != nil {
			panic(err)
		}
	}

	result, err := g.emulator.ExecuteScript(context.Background(), []byte(request.Script), cadenceArgs)

	if err != nil {
		panic(err)
	}

	return string(jsoncdc.MustEncode(result))
}

type ExecuteScriptAtHeightRequest struct {
	ExecuteScriptRequest
	Height uint64 `json:"height"`
}

func (g *InternalGateway) executeScriptAtHeight(this js.Value, args []js.Value) interface{} {
	var request ExecuteScriptAtHeightRequest
	err := json.Unmarshal([]byte(args[0].String()), &request)

	if err != nil {
		panic(err)
	}

	var cadenceArgs []cadence.Value

	if request.ArgsJSON != "" {
		cadenceArgs, err = arguments.ParseJSON(request.ArgsJSON)

		if err != nil {
			panic(err)
		}
	}

	result, err := g.emulator.ExecuteScriptAtHeight(context.Background(), []byte(request.Script), cadenceArgs, request.Height)

	if err != nil {
		panic(err)
	}

	return string(jsoncdc.MustEncode(result))
}

type ExecuteScriptAtID struct {
	ExecuteScriptRequest
	ID string `json:"id"`
}

func (g *InternalGateway) executeScriptAtID(this js.Value, args []js.Value) interface{} {
	var request ExecuteScriptAtID
	err := json.Unmarshal([]byte(args[0].String()), &request)

	if err != nil {
		panic(err)
	}

	var cadenceArgs []cadence.Value

	if request.ArgsJSON != "" {
		cadenceArgs, err = arguments.ParseJSON(request.ArgsJSON)

		if err != nil {
			panic(err)
		}
	}

	result, err := g.emulator.ExecuteScriptAtID(context.Background(), []byte(request.Script), cadenceArgs, sdk.HexToID(request.ID))

	if err != nil {
		panic(err)
	}

	return string(jsoncdc.MustEncode(result))
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
			"collectionId": value.CollectionID.Hex(),
		})
	}

	serializedBlockSeals := make([]interface{}, 0)
	for _, value := range block.Seals {
		serializedBlockSeals = append(serializedBlockSeals, map[string]interface{}{
			"blockId":            value.BlockID.Hex(),
			"executionReceiptId": value.ExecutionReceiptID.Hex(),
		})
	}

	// https://developers.flow.com/tools/clients/fcl-js/api#blockobject
	return map[string]interface{}{
		"id":                   block.ID.Hex(),
		"parentId":             block.ParentID.Hex(),
		"height":               block.Height,
		"timestamp":            block.Timestamp.String(),
		"collectionGuarantees": serializedCollectionGuarantees,
		"blockSeals":           serializedBlockSeals,
		"signatures":           []interface{}{}, // TODO: Implement
	}
}

func (g *InternalGateway) getCollection(this js.Value, args []js.Value) interface{} {
	collection, err := g.emulator.GetCollection(context.Background(), sdk.HexToID(args[0].String()))

	if err != nil {
		panic(err)
	}

	return serializeCollection(collection)
}

func serializeCollection(collection *sdk.Collection) interface{} {
	serializedTransactionIds := make([]interface{}, 0)
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
