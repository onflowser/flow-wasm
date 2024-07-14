# Flow Wasm

Run Flow blockchain emulator on the web.

## Development notes

<details>
    <summary>Invoking the wrapped Go function from JavaScript will pause the event loop and spawn a new goroutine</summary>
Calling a Go function that in turn calls an async JS function will cause a deadlock state.

See: https://withblue.ink/2020/10/03/go-webassembly-http-requests-and-promises.html

</details>
