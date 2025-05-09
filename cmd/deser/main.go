package main

import (
	"fmt"
	"io"
	"os"

	"github.com/dcaravel/stackrox-deserialize/internal/decode"
	"github.com/dcaravel/stackrox-deserialize/internal/encode"
	"github.com/dcaravel/stackrox-deserialize/internal/util"
)

func main() {
	if !isPiped() {
		fmt.Printf("missing data - provide via pipe/stdin\n")
		os.Exit(1)
	}

	dataB, err := io.ReadAll(os.Stdin)
	util.Check(err)

	dataB, err = decode.Hex(dataB)
	util.Check(err)

	entries, err := encode.JSONAll(dataB)
	util.Check(err)

	for _, e := range entries {
		fmt.Printf("====== %s\n", e.Name)
		fmt.Printf("%s\n", e.ProtoJSON)
	}
}

func isPiped() bool {
	fileInfo, _ := os.Stdin.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) == 0
}
