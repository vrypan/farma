package utils

import (
	"encoding/hex"
)

func HexToBytes(str string) []byte {
	if str[0:2] == "0x" {
		str = str[2:]
	}
	bytes, err := hex.DecodeString(str)
	if err != nil {
		panic(err)
	}
	return bytes
}
func BytesToHex(bytes []byte) string {
	return "0x" + hex.EncodeToString(bytes)
}
