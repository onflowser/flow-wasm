cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
CGO_ENABLED=0 GOOS=js GOARCH=wasm go build -tags=no_cgo -o flow-cli.wasm main.go
