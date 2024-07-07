package js

import (
	"context"
	"encoding/json"
	"github.com/onflow/cadence"
	sdk "github.com/onflow/flow-go-sdk"
	"github.com/onflow/flowkit/v2/gateway"
	"syscall/js"
)

// ExternalGateway exposes access API to external networks
// that can't be accessed within WASM (e.g. mainnet, testnet).
type ExternalGateway struct {
	// JS object that implements gateway.Gateway methods
	target js.Value
}

func NewExternalGateway(target js.Value) *ExternalGateway {
	return &ExternalGateway{target: target}
}

func (g *ExternalGateway) GetAccount(ctx context.Context, address sdk.Address) (*sdk.Account, error) {
	value, err := parseResult(resolvePromise(g.target.Call("getAccount", address.Hex())))

	if err != nil {
		return nil, err
	}

	var contracts map[string]string
	err = json.Unmarshal([]byte(value.Get("contracts").String()), &contracts)

	if err != nil {
		panic(err)
	}

	deserializedContracts := make(map[string][]byte)

	for key, value := range contracts {
		deserializedContracts[key] = []byte(value)
	}

	return &sdk.Account{
		Address:   sdk.HexToAddress(value.Get("address").String()),
		Balance:   uint64(value.Get("balance").Int()),
		Contracts: deserializedContracts,
	}, nil
}

func (g *ExternalGateway) SendSignedTransaction(ctx context.Context, transaction *sdk.Transaction) (*sdk.Transaction, error) {
	//TODO implement me
	panic("implement me")
}

func (g *ExternalGateway) GetTransaction(ctx context.Context, identifier sdk.Identifier) (*sdk.Transaction, error) {
	//TODO implement me
	panic("implement me")
}

func (g *ExternalGateway) GetTransactionResultsByBlockID(ctx context.Context, blockID sdk.Identifier) ([]*sdk.TransactionResult, error) {
	//TODO implement me
	panic("implement me")
}

func (g *ExternalGateway) GetTransactionResult(ctx context.Context, identifier sdk.Identifier, b bool) (*sdk.TransactionResult, error) {
	//TODO implement me
	panic("implement me")
}

func (g *ExternalGateway) GetTransactionsByBlockID(ctx context.Context, identifier sdk.Identifier) ([]*sdk.Transaction, error) {
	//TODO implement me
	panic("implement me")
}

func (g *ExternalGateway) ExecuteScript(ctx context.Context, bytes []byte, values []cadence.Value) (cadence.Value, error) {
	//TODO implement me
	panic("implement me")
}

func (g *ExternalGateway) ExecuteScriptAtHeight(ctx context.Context, bytes []byte, values []cadence.Value, u uint64) (cadence.Value, error) {
	//TODO implement me
	panic("implement me")
}

func (g *ExternalGateway) ExecuteScriptAtID(ctx context.Context, bytes []byte, values []cadence.Value, identifier sdk.Identifier) (cadence.Value, error) {
	//TODO implement me
	panic("implement me")
}

func (g *ExternalGateway) GetLatestBlock(ctx context.Context) (*sdk.Block, error) {
	//TODO implement me
	panic("implement me")
}

func (g *ExternalGateway) GetBlockByHeight(ctx context.Context, u uint64) (*sdk.Block, error) {
	//TODO implement me
	panic("implement me")
}

func (g *ExternalGateway) GetBlockByID(ctx context.Context, identifier sdk.Identifier) (*sdk.Block, error) {
	//TODO implement me
	panic("implement me")
}

func (g *ExternalGateway) GetEvents(ctx context.Context, s string, u uint64, u2 uint64) ([]sdk.BlockEvents, error) {
	//TODO implement me
	panic("implement me")
}

func (g *ExternalGateway) GetCollection(ctx context.Context, identifier sdk.Identifier) (*sdk.Collection, error) {
	//TODO implement me
	panic("implement me")
}

func (g *ExternalGateway) GetLatestProtocolStateSnapshot(ctx context.Context) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (g *ExternalGateway) Ping() error {
	//TODO implement me
	panic("implement me")
}

func (g *ExternalGateway) WaitServer(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (g *ExternalGateway) SecureConnection() bool {
	//TODO implement me
	panic("implement me")
}

var _ gateway.Gateway = &ExternalGateway{}
