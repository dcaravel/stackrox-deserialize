package decode

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

func Hex(dataB []byte) ([]byte, error) {
	dataB = bytes.TrimSpace(dataB)

	for dataB[0] == '\\' || dataB[0] == 'x' {
		dataB = dataB[1:]
	}

	cleanedDataB := []byte{}
	for _, b := range dataB {
		if b != '\n' {
			cleanedDataB = append(cleanedDataB, b)
		}
	}
	dataB = cleanedDataB

	protoBytes, err := hex.DecodeString(string(dataB))
	if err != nil {
		return nil, fmt.Errorf("decoding string to hex: %w", err)
	}

	return protoBytes, nil
}
