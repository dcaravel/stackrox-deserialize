# StackRox Deserialize
A static site that can deserialized hex encoded StackRox proto messages to assist troubleshooting.

https://dcaravel.github.io/stackrox-deserialize/

The deserialization happens via a Go built `wasm` module.

## Build

```
GOOS=js GOARCH=wasm go build -o docs/main.wasm ./cmd/wasm
```

## Run
Can use any web server, an example:

```
cd <somewhere else>
git clone git@github.com:davidwashere/daserve.git
cd daserve
go install .

cd <this repo>
daserve docs
```
