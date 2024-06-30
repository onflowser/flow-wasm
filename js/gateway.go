package js

import (
	"context"
	"github.com/onflow/cadence"
	sdk "github.com/onflow/flow-go-sdk"
	"github.com/onflow/flowkit/v2/gateway"
	"syscall/js"
)

type Gateway struct {
	// JS object that implements Gateway methods
	target js.Value
}

func NewGateway(target js.Value) *Gateway {
	return &Gateway{target: target}
}

func (g *Gateway) GetAccount(ctx context.Context, address sdk.Address) (*sdk.Account, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Gateway) SendSignedTransaction(ctx context.Context, transaction *sdk.Transaction) (*sdk.Transaction, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Gateway) GetTransaction(ctx context.Context, identifier sdk.Identifier) (*sdk.Transaction, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Gateway) GetTransactionResultsByBlockID(ctx context.Context, blockID sdk.Identifier) ([]*sdk.TransactionResult, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Gateway) GetTransactionResult(ctx context.Context, identifier sdk.Identifier, b bool) (*sdk.TransactionResult, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Gateway) GetTransactionsByBlockID(ctx context.Context, identifier sdk.Identifier) ([]*sdk.Transaction, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Gateway) ExecuteScript(ctx context.Context, bytes []byte, values []cadence.Value) (cadence.Value, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Gateway) ExecuteScriptAtHeight(ctx context.Context, bytes []byte, values []cadence.Value, u uint64) (cadence.Value, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Gateway) ExecuteScriptAtID(ctx context.Context, bytes []byte, values []cadence.Value, identifier sdk.Identifier) (cadence.Value, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Gateway) GetLatestBlock(ctx context.Context) (*sdk.Block, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Gateway) GetBlockByHeight(ctx context.Context, u uint64) (*sdk.Block, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Gateway) GetBlockByID(ctx context.Context, identifier sdk.Identifier) (*sdk.Block, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Gateway) GetEvents(ctx context.Context, s string, u uint64, u2 uint64) ([]sdk.BlockEvents, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Gateway) GetCollection(ctx context.Context, identifier sdk.Identifier) (*sdk.Collection, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Gateway) GetLatestProtocolStateSnapshot(ctx context.Context) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Gateway) Ping() error {
	//TODO implement me
	panic("implement me")
}

func (g *Gateway) WaitServer(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (g *Gateway) SecureConnection() bool {
	//TODO implement me
	panic("implement me")
}

var _ gateway.Gateway = &Gateway{}
