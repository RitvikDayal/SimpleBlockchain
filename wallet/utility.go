package wallet

import (
	"github.com/mr-tron/base58"
)

func Base58Encode(data []byte) []byte {
	encode := base58.Encode(data)

	return []byte(encode)
}

func Base58Decode(data []byte) []byte {
	decode, err := base58.Decode(string(data))
	if err != nil {
		panic(err)
	}

	return decode
}