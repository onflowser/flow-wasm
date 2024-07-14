package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/onflow/flow-emulator/storage/memstore"
	"github.com/onflow/flowkit/v2"
	"github.com/onflow/flowkit/v2/config"
	"github.com/onflow/flowkit/v2/deps"
	"github.com/onflow/flowkit/v2/output"
	jsFlow "github.com/onflowser/flow-cli-wasm/js"
	"github.com/onflowser/flow-cli-wasm/logging"
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
	logger    *logging.Logger
	kit       *flowkit.Flowkit
	installer *deps.DependencyInstaller
}

func main() {
	w := New(Config{
		Verbose:    true,
		LogFormat:  "text",
		Prompter:   jsFlow.NewPrompter(js.Global().Get("prompter")),
		FileSystem: jsFlow.NewFileSystem(js.Global().Get("flowFileSystem")),
	})

	// Register APIs
	internalGateway := jsFlow.NewInternalGateway(w.gateway)
	js.Global().Set("gateway", internalGateway.JsValue())
	js.Global().Set("getLogs", js.FuncOf(w.getLogs))
	js.Global().Set("install", js.FuncOf(w.install))
	js.Global().Set("deploy", js.FuncOf(w.deploy))

	// Indicate the emulator started and APIs were initialized
	js.Global().Call("onStarted")

	// Prevent the function from returning, which is required in a wasm module
	select {}
}

func New(config Config) *FlowWasm {
	logger := logging.NewLogger(logging.Config{
		Verbose:   config.Verbose,
		LogFormat: config.LogFormat,
	})
	store := memstore.New()

	emulatorGateway := gateway.NewEmulatorGatewayWithOpts(
		&gateway.EmulatorKey{
			PublicKey: emulator.DefaultServiceKey().AccountKey().PublicKey,
			SigAlgo:   emulator.DefaultServiceKeySigAlgo,
			HashAlgo:  emulator.DefaultServiceKeyHashAlgo,
		},
		gateway.WithEmulatorOptions(
			emulator.WithLogger(*logger.Zerolog()),
			emulator.WithStore(store),
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

	kit := flowkit.NewFlowkit(state, *network, emulatorGateway, logger)

	installer, err := deps.NewDependencyInstaller(
		state,
		config.Prompter,
		deps.WithGateways(jsGateways(emulatorGateway)),
		deps.WithLogger(logger),
		deps.WithSaveState(),
	)

	if err != nil {
		panic(err)
	}

	return &FlowWasm{
		config:    config,
		gateway:   emulatorGateway,
		logger:    logger,
		kit:       kit,
		installer: installer,
	}
}

func jsGateways(emulatorGateway gateway.Gateway) map[string]gateway.Gateway {
	testnetGateway := jsFlow.NewExternalGateway(js.Global().Get("testnetGateway"))
	mainnetGateway := jsFlow.NewExternalGateway(js.Global().Get("mainnetGateway"))
	previewnetGateway := jsFlow.NewExternalGateway(js.Global().Get("previewnetGateway"))

	return map[string]gateway.Gateway{
		config.EmulatorNetwork.Name:   emulatorGateway,
		config.TestnetNetwork.Name:    testnetGateway,
		config.MainnetNetwork.Name:    mainnetGateway,
		config.PreviewnetNetwork.Name: previewnetGateway,
	}
}

func (w *FlowWasm) install(this js.Value, args []js.Value) any {
	executor := func() (js.Value, error) {
		err := w.installer.Install()
		return js.Null(), err
	}

	return jsFlow.AsyncWork(executor)
}

func (w *FlowWasm) getLogs(this js.Value, args []js.Value) interface{} {
	res, err := json.Marshal(w.logger.LogsHistory())

	if err != nil {
		panic(err)
	}

	return string(res)
}

func (w *FlowWasm) deploy(this js.Value, args []js.Value) interface{} {
	executor := func() (js.Value, error) {
		contracts, err := w.kit.DeployProject(context.Background(), flowkit.UpdateExistingContract(true))
		if err != nil {
			var projectErr *flowkit.ProjectDeploymentError
			if errors.As(err, &projectErr) {
				for name, err := range projectErr.Contracts() {
					w.logger.Info(fmt.Sprintf(
						"%s Failed to deploy contract %s: %s",
						output.ErrorEmoji(),
						name,
						err.Error(),
					))
				}
				return js.Null(), fmt.Errorf("failed deploying all contracts")
			}

			return js.Null(), err
		}

		for _, contract := range contracts {
			w.logger.Info(fmt.Sprintf("deployed %s contract", contract.Name))
		}

		return js.Null(), err
	}

	return jsFlow.AsyncWork(executor)
}
