package main

import (
	"fmt"
)

// LinkedList Structure that will be used to store blocks

// Block will be chained up in a linked list
// include helper methods that allow easy access to find parentBlockHash, hash, and nonce
type Block struct {
	header BlockHeader
	data   []Transaction // transactions need to be in a linked list
}

type BlockHeader struct {
	timestamp       int64
	parentBlockHash []byte // data type used in ethash library
	hash            []byte // block seed takes in uint64
	nonce           uint64 // rand int
}

type Transaction struct {
	sender    []byte
	recipient []byte
	timestamp int64
	signature []byte // from MerkleTree library
	data      []byte
}

func main() {
	fmt.Println("Test.")
}
