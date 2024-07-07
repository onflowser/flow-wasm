package js

import (
	"context"
	"github.com/onflow/cadence"
	sdk "github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flowkit/v2/gateway"
	"syscall/js"
)

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

	serializedContracts := make(map[string]interface{})

	for key, value := range account.Contracts {
		serializedContracts[key] = string(value)
	}

	serializedKeys := make([]interface{}, 0)

	// Mapping as defined in https://github.com/onflow/fcl-js/blob/9c7873140015c9d1e28712aed93c56654f656639/packages/transport-grpc/src/send-get-account.js#L16-L28
	hashAlgoToJsIndex := map[crypto.HashAlgorithm]int{
		crypto.SHA2_256: 1,
		crypto.SHA2_384: 2,
		crypto.SHA3_256: 3,
		crypto.SHA3_384: 4,
		crypto.KMAC128:  5,
	}

	// Mapping as defined in https://github.com/onflow/fcl-js/blob/9c7873140015c9d1e28712aed93c56654f656639/packages/transport-grpc/src/send-get-account.js#L16-L28
	signAlgoToJsIndex := map[crypto.SignatureAlgorithm]int{
		crypto.ECDSA_P256:      1,
		crypto.ECDSA_secp256k1: 2,
		crypto.BLS_BLS12_381:   3,
	}

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

	return map[string]interface{}{
		"address":   account.Address.String(),
		"balance":   account.Balance,
		"contracts": serializedContracts,
		"keys":      serializedKeys,
		"code":      string(account.Code),
	}
}

func (g *InternalGateway) sendSignedTransaction(ctx context.Context, transaction *sdk.Transaction) (*sdk.Transaction, error) {
	//TODO implement me
	panic("implement me")
}

func (g *InternalGateway) getTransaction(ctx context.Context, identifier sdk.Identifier) (*sdk.Transaction, error) {
	//TODO implement me
	panic("implement me")
}

func (g *InternalGateway) getTransactionResultsByBlockID(ctx context.Context, blockID sdk.Identifier) ([]*sdk.TransactionResult, error) {
	//TODO implement me
	panic("implement me")
}

func (g *InternalGateway) getTransactionResult(ctx context.Context, identifier sdk.Identifier, b bool) (*sdk.TransactionResult, error) {
	//TODO implement me
	panic("implement me")
}

func (g *InternalGateway) getTransactionsByBlockID(ctx context.Context, identifier sdk.Identifier) ([]*sdk.Transaction, error) {
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

func (g *InternalGateway) getLatestBlock(ctx context.Context) (*sdk.Block, error) {
	//TODO implement me
	panic("implement me")
}

func (g *InternalGateway) getBlockByHeight(ctx context.Context, u uint64) (*sdk.Block, error) {
	//TODO implement me
	panic("implement me")
}

func (g *InternalGateway) getBlockByID(ctx context.Context, identifier sdk.Identifier) (*sdk.Block, error) {
	//TODO implement me
	panic("implement me")
}

func (g *InternalGateway) getEvents(ctx context.Context, s string, u uint64, u2 uint64) ([]sdk.BlockEvents, error) {
	//TODO implement me
	panic("implement me")
}

func (g *InternalGateway) getCollection(ctx context.Context, identifier sdk.Identifier) (*sdk.Collection, error) {
	//TODO implement me
	panic("implement me")
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
