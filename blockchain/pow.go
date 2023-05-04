package blockchain

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
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
// Lsh() is a left shift, which is a bitwise operation
// sets z = x << n and returns z
// z = 1 * 2^T(256 - difficulty)
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
func (block *Block) BlockDataToBytes() []byte {
	data := bytes.Join(
		[][]byte{
			block.header.parentBlockHash,
			block.GetData().Root.Tree.MerkleRoot(), // root hash of merkle tree
			ToHex(int64(block.GetNonce())),
			ToHex(int64(difficulty)),
		},
		[]byte{},
	)
	return data
}

// Mining of the block.
// Gets a specific hash that is less than the target hash
// Returns nonce and hash
func (block *Block) Mine() (int, [32]byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0

	for nonce < math.MaxInt64 {
		hash, _ := block.CalculateHash()

		// prints out all the hashes
		/*
			fmt.Printf("\r%x", hash) // \r is carriage return, %x is hex
			fmt.Println()
		*/
		intHash.SetBytes(hash[:])

		// if intHash < target, then we have found a valid hash
		// Cmp compares x and y and returns:
		// -1 if x < y
		// 0 if x == y
		// +1 if x > y
		if intHash.Cmp(block.GetTarget()) == -1 {
			//fmt.Printf("%x", hash)
			block.SetHash(hash)
			fmt.Printf("Hash has been set: %x", hash)
			break
		} else {
			nonce++
			block.SetNonce(uint64(nonce))
		}
	}
	fmt.Println()

	return nonce, hash
}

func (block *Block) TestPOW(newDiff int) {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-newDiff)) // left shift

	pow := &ProofOfWork{target}

	block.pow = pow
}

func (block *Block) TestPrintMine() (int, [32]byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0

	for nonce < math.MaxInt64 {
		hash, _ := block.CalculateHash()

		// prints out all the hashes
		fmt.Printf("\r%x", hash) // \r is carriage return, %x is hex
		fmt.Println()

		intHash.SetBytes(hash[:])

		if intHash.Cmp(block.GetTarget()) == -1 {
			block.SetHash(hash)
			fmt.Printf("Hash has been set: %x", hash)
			break
		} else {
			nonce++
			block.SetNonce(uint64(nonce))
		}
	}
	fmt.Println()

	return nonce, hash
}

/*
// Might not need to use, have "VerifyContent" in MerkleTree package
func (block *Block) Validate() bool {
    var intHash big.Int
	hash, _ := block.CalculateHash()

    intHash.SetBytes(hash[:])

    return intHash.Cmp(block.GetTarget()) == -1
}
*/
