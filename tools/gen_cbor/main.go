package main

import (
	"github.com/ipfs-force-community/venus-messager/chain/types"
	gen "github.com/whyrusleeping/cbor-gen"
)

func main() {
	if err := gen.WriteTupleEncodersToFile("./chain/types/cbor_gen.go", "types",
		types.Message{},
		types.SignedMessage{},
	); err != nil {
		panic(err)
	}
}
