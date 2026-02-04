# StackRox Deserialize
A static site that can deserialized hex encoded StackRox proto messages to assist troubleshooting.

https://dcaravel.github.io/stackrox-deserialize/

The deserialization happens via a Go built `wasm` module.

## Updating stackrox version

```sh
# expected to fail, but copy the version string from the failure
go get github.com/stackrox/stackrox@latest

# then tidy so that go.sum is also updated
go mod tidy

```

## Build

```
GOOS=js GOARCH=wasm go build -o docs/main.wasm ./cmd/wasm
```
This 'build' when pushed will cause GH pages to be updated

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


## Other

Updating [monaco editor](https://www.npmjs.com/package/monaco-editor)
```sh
npm install monaco-editor
...
cp -r node_modules/monaco-editor/min docs/monaco-editor/min
```