//go:build js

package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
	"time"

	"github.com/dcaravel/stackrox-deserialize/internal/decode"
	"github.com/dcaravel/stackrox-deserialize/internal/encode"
)

type Out struct {
	Name string    `json:"name,omitempty"`
	Date time.Time `json:"date,omitempty"`
}

func main() {
	fmt.Println("Go wasm module started")
	c := make(chan struct{}, 0)

	registerJSFuncs()
	fmt.Println("Go funcs registered")

	<-c
	fmt.Println("Go wasm module exited")
}

func registerJSFuncs() {
	js.Global().Set("decode", FuncPromise(Decode))
}

func Decode(args []js.Value) (interface{}, error) {
	data := args[0].String()
	dataB := []byte(data)

	dataB, err := decode.Hex(dataB)
	if err != nil {
		return nil, fmt.Errorf("decoding value into hex: %w", err)
	}

	entries, err := encode.JSONAll(dataB)
	if err != nil {
		return nil, fmt.Errorf("converting value to json: %w", err)
	}

	return entries, nil
}

// FuncPromise wraps a Go function that returns (result, error) into a JS Promise-returning function
func FuncPromise(fn func(args []js.Value) (interface{}, error)) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		promiseConstructor := js.Global().Get("Promise")

		// Create the executor function for the Promise
		executor := js.FuncOf(func(_ js.Value, executorArgs []js.Value) interface{} {
			resolve := executorArgs[0]
			reject := executorArgs[1]

			go func() {
				defer func() {
					if r := recover(); r != nil {
						reject.Invoke(fmt.Sprintf("panic: %v", r))
					}
				}()

				result, err := fn(args)
				if err != nil {
					reject.Invoke(err.Error())
				} else {
					dataB, err := json.MarshalIndent(result, "", "  ")
					if err != nil {
						reject.Invoke(err.Error())
					}

					resolve.Invoke(string(dataB))
				}
			}()

			return nil
		})

		return promiseConstructor.New(executor)
	})
}
