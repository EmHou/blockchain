package blockchain

import (
	"math/big"
	"encoding/binary"
	"bytes"
	"log"
)

// Inspired by Noah Hein's "Building a Blockchain in Go PT:II - Proof of Work"
// https://dev.to/nheindev/building-a-blockchain-in-go-pt-ii-proof-of-work-eel
// Edited for our implementation

// difficulty in our implementation is set to 12 for simplification
// In the real world, difficulty is set in each block and is based on the number of nodes in the network
const difficulty = 12 

type ProofOfWork struct {
	target *big.Int
}

// Hein states: The closer we get to 256, the easier the computation will be. Increasing our difficulty will increase the runtime of our algorithm
func NewPOW(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty)) // left shift

	pow := &ProofOfWork{target}

	return pow
}

// ToHex converts int64 to []byte
func ToHex(num int64) []byte {
    buff := new(bytes.Buffer)
    err := binary.Write(buff, binary.BigEndian, num)
    if err != nil {
        log.Panic(err)
    }
    return buff.Bytes()
}

// Takes nonce, timestamps, parentBlockHash, and root hash of the transaction Merkle Tree
// Returns []byte for Nonce
func (block *Block) InitNonce(nonce int) []byte {
    data := bytes.Join(
        [][]byte{
            block.header.parentBlockHash,
            block.getData().Root.Tree.MerkleRoot(), // root hash of merkle tree
            ToHex(int64(nonce)),
            ToHex(int64(difficulty)),
        },
        []byte{},
    )
    return data
}




