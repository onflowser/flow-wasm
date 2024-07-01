package js

import (
	"fmt"
	"syscall/js"
)

// AsyncWork handles executing blocking code in go routine with promises to avoid deadlocks
// See: https://withblue.ink/2020/10/03/go-webassembly-http-requests-and-promises.html
func AsyncWork(executor func() (js.Value, error)) js.Value {
	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {
			result, err := executor()
			if err == nil {
				resolve.Invoke(result)
			} else {
				reject.Invoke(err.Error())
			}
		}()

		// The handler of a Promise doesn't return any value
		return nil
	})

	// Create and return the Promise object
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// resolvePromise awaits any JS promise-like value with a "then" function (aka. thenable)
// Errors passed to "catch" will be ignored.
// Instead, the promise should return a GoResult object with an "error" property.
func resolvePromise(promise js.Value) js.Value {
	var result js.Value
	wait := make(chan interface{})
	go func() {
		promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			result = args[0]
			wait <- nil
			return nil
		}))
	}()
	<-wait
	return result
}

// parseResult accepts a JS object that implements GoResult interface from /ts-lib/src/go-interfaces.ts
func parseResult(jsObject js.Value) (js.Value, error) {
	value := jsObject.Get("value")
	rawErr := jsObject.Get("error")

	if !rawErr.IsNull() {
		return value, fmt.Errorf(rawErr.String())
	}

	return value, nil
}
