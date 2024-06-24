package main

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"os"
	"syscall/js"

	"github.com/onflow/flow-emulator/emulator"
	flowgo "github.com/onflow/flow-go/model/flow"
)

type Config struct {
	Verbose   bool
	LogFormat string // "text" or "json". Defaults to "json" if "logs" writer is used.
}

type WasmEmulator struct {
	config     Config
	blockchain *emulator.Blockchain
	logger     *zerolog.Logger
	cachedLogs *CacheLogWriter
}

func main() {
	fmt.Println("Starting emulator")

	w := NewWasmEmulator(Config{
		Verbose:   true,
		LogFormat: "text",
	})

	// Mount the function on the JavaScript global object.
	js.Global().Set("GetAccount", js.FuncOf(w.GetAccount))
	js.Global().Set("GetLogs", js.FuncOf(w.GetLogs))

	// Prevent the function from returning, which is required in a wasm module
	select {}
}

func NewWasmEmulator(config Config) *WasmEmulator {
	logger, cacheWriter := initLogger(config)

	blockchain, err := emulator.New(
		emulator.WithLogger(*logger),
	)

	if err != nil {
		panic(err)
	}

	return &WasmEmulator{
		config:     config,
		blockchain: blockchain,
		logger:     logger,
		cachedLogs: cacheWriter,
	}
}

func (w *WasmEmulator) GetAccount(this js.Value, args []js.Value) interface{} {
	address := flowgo.HexToAddress(args[0].String())
	account, err := w.blockchain.GetAccount(address)

	if err != nil {
		panic(err)
	}

	return map[string]interface{}{
		"address": account.Address.String(),
		"balance": account.Balance,
		// "contracts": account.Contracts,
	}
}

func (w *WasmEmulator) GetLogs(this js.Value, args []js.Value) interface{} {
	fmt.Println("Cache size", len(w.cachedLogs.logs))

	res, err := json.Marshal(&w.cachedLogs.logs)

	if err != nil {
		panic(err)
	}

	return string(res)
}

func initLogger(config Config) (*zerolog.Logger, *CacheLogWriter) {

	level := zerolog.InfoLevel
	if config.Verbose {
		level = zerolog.DebugLevel
	}
	zerolog.MessageFieldName = "msg"

	cacheWriter := NewCacheLogWriter()

	writer := zerolog.MultiLevelWriter(
		NewTextWriter(),
		cacheWriter,
	)

	logger := zerolog.New(writer).With().Timestamp().Logger().Level(level)

	return &logger, cacheWriter
}

func NewTextWriter() zerolog.ConsoleWriter {
	writer := zerolog.ConsoleWriter{Out: os.Stdout}
	writer.FormatMessage = func(i interface{}) string {
		if i == nil {
			return ""
		}
		return fmt.Sprintf("%-44s", i)
	}

	return writer
}

type CacheLogWriter struct {
	logs []string
}

func NewCacheLogWriter() *CacheLogWriter {
	return &CacheLogWriter{
		logs: make([]string, 0),
	}
}

var _ io.Writer = &CacheLogWriter{}

func (c *CacheLogWriter) Write(p []byte) (n int, err error) {
	c.logs = append(c.logs, string(p))
	return len(p), nil
}
