package main

import (
	"context"
	"encoding/json"
	"github.com/onflow/flow-emulator/storage/memstore"
	sdk "github.com/onflow/flow-go-sdk"
	"github.com/onflow/flowkit/v2"
	"github.com/onflow/flowkit/v2/deps"
	jsFlow "github.com/onflowser/flow-cli-wasm/js"
	"github.com/onflowser/flow-cli-wasm/logger"
	"syscall/js"

	"github.com/onflow/flow-emulator/emulator"
	"github.com/onflow/flowkit/v2/gateway"
)

type Config struct {
	Verbose    bool
	LogFormat  string // "text" or "json". Defaults to "json" if "logs" writer is used.
	FileSystem flowkit.ReaderWriter
	Prompter   deps.Prompter
}

type FlowWasm struct {
	config    Config
	state     *flowkit.State
	gateway   *gateway.EmulatorGateway
	logger    *logger.Logger
	kit       *flowkit.Flowkit
	installer *deps.DependencyInstaller
}

func main() {
	w := New(Config{
		Verbose:    true,
		LogFormat:  "text",
		FileSystem: jsFlow.NewFileSystem(js.Global().Get("flowFileSystem")),
	})

	// Mount the function on the JavaScript global object.
	js.Global().Set("GetAccount", js.FuncOf(w.GetAccount))
	js.Global().Set("GetLogs", js.FuncOf(w.GetLogs))

	// Prevent the function from returning, which is required in a wasm module
	select {}
}

func New(config Config) *FlowWasm {
	l := logger.NewLogger(logger.Config{
		Verbose:   config.Verbose,
		LogFormat: config.LogFormat,
	})
	s := memstore.New()

	g := gateway.NewEmulatorGatewayWithOpts(
		&gateway.EmulatorKey{
			PublicKey: emulator.DefaultServiceKey().AccountKey().PublicKey,
			SigAlgo:   emulator.DefaultServiceKeySigAlgo,
			HashAlgo:  emulator.DefaultServiceKeyHashAlgo,
		},
		gateway.WithEmulatorOptions(
			emulator.WithLogger(*l.Zerolog()),
			emulator.WithStore(s),
			emulator.WithTransactionValidationEnabled(false),
			emulator.WithStorageLimitEnabled(false),
			emulator.WithTransactionFeesEnabled(false),
		),
	)

	configFilePaths := []string{
		"flow.json",
	}
	state, err := flowkit.Load(configFilePaths, config.FileSystem)
	if err != nil {
		panic(err)
	}

	network, err := state.Networks().ByName("emulator")
	if err != nil {
		panic(err)
	}

	kit := flowkit.NewFlowkit(state, *network, g, l)

	installer, err := deps.NewDependencyInstaller(state, config.Prompter)

	if err != nil {
		panic(err)
	}

	return &FlowWasm{
		config:    config,
		gateway:   g,
		logger:    l,
		kit:       kit,
		installer: installer,
	}
}

func (w *FlowWasm) GetAccount(this js.Value, args []js.Value) interface{} {
	account, err := w.gateway.GetAccount(context.Background(), sdk.HexToAddress(args[0].String()))

	if err != nil {
		panic(err)
	}

	return map[string]interface{}{
		"address": account.Address.String(),
		"balance": account.Balance,
		// "contracts": account.Contracts,
	}
}

func (w *FlowWasm) GetLogs(this js.Value, args []js.Value) interface{} {
	res, err := json.Marshal(w.logger.LogsHistory())

	if err != nil {
		panic(err)
	}

	return string(res)
}
