if [ ! -f ./wasm_exec.js ]; then
    # This should work, but doesn't: https://github.com/actions/setup-go/issues/71
    # cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
    curl https://raw.githubusercontent.com/golang/go/go1.22.3/misc/wasm/wasm_exec.js > wasm_exec.js
fi

CGO_ENABLED=0 GOOS=js GOARCH=wasm go build -tags=no_cgo -o flow.wasm main.go
