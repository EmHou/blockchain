package blockchain

import (
	"fmt"
	"bytes"
	"encoding/binary"
	"log"
	"math/big"
	"math"
	"crypto/sha256"
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
func NewPOW() *ProofOfWork {
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
			block.GetData().Root.Tree.MerkleRoot(), // root hash of merkle tree
			ToHex(int64(nonce)),
			ToHex(int64(difficulty)),
		},
		[]byte{},
	)
	return data
}


// reimpliment this using Merkle Tree interface
// need to implement CalculateHash()
func (block *Block) Run() (int, []byte) {
	var intHash big.Int
	var hash []byte

	nonce := 0

	for nonce < math.MaxInt64 {
		hash, _ = block.CalculateHash(nonce)

		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:])

		// if intHash < target, then we have found a valid hash
		// Cmp compares x and y and returns: 
		// -1 if x < y
		// 0 if x == y
		// +1 if x > y
		if intHash.Cmp(block.GetTarget()) == -1 {
			break
		} else {
			nonce++
		}
	}

	return nonce, hash[:] /// [:] makes a slice of the array
}

func (block *Block) Validate() bool {
    var intHash big.Int

    data := block.InitNonce(int(block.GetNonce()))

    hash := sha256.Sum256(data)
    intHash.SetBytes(hash[:])

    return intHash.Cmp(block.GetTarget()) == -1
}
