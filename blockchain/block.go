package blockchain

import (
	"math/big"
	"time"
	"crypto/sha256"

	"github.com/cbergoon/merkletree"
)

// BlockChain structure links together blocks in a Merkle Tree
// currentBlock		block that is being filled up with transactions has NOT been added to the chain yet
// chain			chain of blocks that have been added to the chain
type BlockChain struct {
	currentBlock *Block
	chain *merkletree.MerkleTree
}

// include helper methods that allow easy access to find parentBlockHash, hash, and nonce
// arbitrary number of transactions - in this implemenation, we choose 7 transactions per block
// header	contains metadata of the block
// data		transactions (Merkle Tree)
type Block struct {
	header BlockHeader
	data   *merkletree.MerkleTree
	pow    *ProofOfWork
}

// timestamp 		time when block is added to the chain
// parentBlockHash	hash of the previous block
// hash	 			takes nonce, timestamps, parentBlockHash, and root hash of the transaction Merkle Tree
// nonce	 		rand int that is initialised to 0
type BlockHeader struct {
	timestamp       int64 
	parentBlockHash []byte
	hash            []byte 
	nonce           uint64 
}

// sender and recipient in cryptocurrency are referring to the keys of the wallets
// in our implementation, we will use the ip address of the sender and recipient. This is only used for metadata for hashing purposes
// not actually sending anything to the recipient (no actual cryptocurrency)
// in real-life cryptocurrency, recipient does not have to be online to receive
// sender		ip address of the sender
// recipient	ip address of the recipient
// timestamp	time when transaction is created
// data			message that sender wants to send to recipient. In real-life cryptocurrency, 
//				this is the amount of cryptocurrency that the sender wants to send to the recipient
type Transaction struct {
	sender    []byte
	recipient []byte 
	timestamp int64
	data      []byte
}

func (Block Block) CalculateHash(nonce int) ([]byte, error) {
	hash := sha256.New()
	data := Block.InitNonce(nonce)

	if _, err := hash.Write(data); err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}

//Equals tests for equality of two Contents
func (block Block) Equals(other merkletree.Content) (bool, error) {
	return block.data == other.(Block).data, nil
}


func (block *Block) GetTimestamp() int64 {
	return block.header.timestamp
}

func (block *Block) GetHeader() BlockHeader {
	return block.header
}

func (block *Block) GetParentBlockHash() []byte {
	return block.header.parentBlockHash
}

func (block *Block) GetHash() []byte {
	return block.header.hash
}

func (block *Block) GetNonce() uint64 {
	return block.header.nonce
}

func (block *Block) GetData() *merkletree.MerkleTree {
	return block.data
}

func (block *Block) GetTarget() *big.Int {
	return block.pow.target
}

// NewBlockChain creates a new block
//
func MakeBlock(pBlockHash []byte) *Block {
	header := &BlockHeader {
		timestamp: time.Now().UnixNano(),
		parentBlockHash: pBlockHash,
		hash: []byte{}, // empty hash that will be set 
		nonce: 0,
	}

	block := &Block {
		header: *header,
		//data: merkletree.NewTree(),
		pow: NewPOW(),
	}

	return block
}





//-- Notes --//

// Step 1:
// Transaction -> get hash -> put in Merkle Tree -> get root hash -> put in block
// Leaves are added in the order that they are hashed (sequential order)
// Note: Merkle Tree is a binary tree, a new node is created with the old root as its left child and the new hash as its right child.
// Step 2:
// Blocks -> get hash -> put in Merkle tree -> get root hash -> put in chain

// Can only add block to chain if it is accepted by the network

