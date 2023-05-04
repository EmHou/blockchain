package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"time"

	"github.com/cbergoon/merkletree"
)

var max int = 7

// Includes helper methods that allow easy access to find parentBlockHash, hash, and nonce.
// Arbitrary number of transactions - in this implemenation, we choose 7 transactions per block.
// header	contains metadata of the block
// data		transactions (Merkle Tree)
// pow		Proof of Work algorithm to validate blocks and calculate nonce to add block to chain
// dataList	transactions that are to added to the block
type Block struct {
	header   BlockHeader
	data     *merkletree.MerkleTree
	pow      *ProofOfWork
	dataList []merkletree.Content
}

// BlockHeader contains metadata of the block.
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

func (block *Block) SetHash(hash []byte) {
	block.header.hash = hash
}

func (block *Block) GetNonce() uint64 {
	return block.header.nonce
}

func (block *Block) SetNonce(nonce uint64) {
	block.header.nonce = nonce
}

func (block *Block) GetData() *merkletree.MerkleTree {
	return block.data
}

func (block *Block) GetTarget() *big.Int {
	return block.pow.target
}

func (block *Block) GetDataList() []merkletree.Content {
	return block.dataList
}

func SetMax(maxiumum int) {
	max = maxiumum
}

func GetMax() int {
	return max
}

// Part of the Content interface in MerkleTree Package.
// TODO
func (Block Block) CalculateHash() ([]byte, error) {
	hash := sha256.New()
	data := Block.BlockDataToBytes()

	if _, err := hash.Write(data); err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}

// need to set hash of the block after it is added to the chain
// can only use this method after the block is added to the chain
// because the data (MerkleTree of transacitons) must be set before the hash can be set
// Part of the Content interface in MerkleTree Package
// TODO
func (block Block) Equals(other merkletree.Content) (bool, error) {
	return block.data == other.(Block).data, nil
}

// NewBlockChain creates a new block and returns the pointer to it.
// Initialises:
// timestamp		to current time in nanoseconds
// parentBlockHash	to the hash of the previous block (gotten from parameter)
// hash				to empty byte array, to be set after a nonce is found (after mining when adding to chain)
// nonce 			to 0
// dataList 		to empty array of type (merkletree.Content)
// pow 				to a new ProofOfWork struct
func MakeBlock(pBlockHash []byte) *Block {
	header := &BlockHeader{
		timestamp:       time.Now().UnixNano(),
		parentBlockHash: pBlockHash,
		hash:            []byte{},
		nonce:           0,
	}

	block := &Block{
		header:   *header,
		dataList: []merkletree.Content{},
		pow:      NewPOW(),
	}

	return block
}

// Takes MakeBlock() and passes empty byte array to it
// Has no parent block, so parentBlockHash is empty (empty byte array)
func MakeGenesisBlock() *Block {
	genesis := MakeBlock([]byte{})

	for i := 0; i < max; i++ {
		emptyTransaction := Transaction{
			Sender:    []byte{},
			Recipient: []byte{},
			Timestamp: time.Now().UnixNano(),
			Data:      []byte("init"),
		}

		genesis.AddTransaction(emptyTransaction)
	}
	genesis.Mine()

	return genesis
}

// AddTransaction is of type Block and takes paramater of type Transaction.
// Checks if the block is full (var max transactions, arbitrarily set) and if it is, prints error message.
// If the transaction MerkleTree is empty, creates a new Merkle Tree with the transaction.
// If it is not empty, rebuilds the Merkle Tree with the new transaction.
func (block *Block) AddTransaction(transaction Transaction) {
	if len(block.dataList) >= max {
		fmt.Println("Block is full, cannot add more transactions")

	} else {
		block.dataList = append(block.dataList, transaction)

		if block.data == nil {
			tree, err := merkletree.NewTree(block.dataList)
			if err != nil {
				log.Fatal(err)
			}

			block.data = tree
		} else {
			err := block.data.RebuildTreeWith(block.dataList)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

}

// Function used to set parent hash
// Used when first block is added to the blockchain after the genesis block.
// Otherwise, the parent hash will be empty for this block, and it will not be added.
func (block *Block) SetBlockParentHash(Hash []byte) {
	block.header.parentBlockHash = Hash
}

func (block Block) String() string {
	str := "**Block**\n"
	str += block.header.String()
	str += "Data (String representation of Transactions Merkle Tree):\n"
	str += block.data.String() + "\n"

	return str
}

func (header *BlockHeader) String() string {
	str := "Timestamp: " + strconv.FormatInt(header.timestamp, 10) + "\n"
	str += "Parent Block Hash: " + hex.EncodeToString(header.parentBlockHash) + "\n"
	str += "Hash: " + hex.EncodeToString(header.hash) + "\n"

	return str
}

//-- Notes --//

// Step 1:
// Transaction -> get hash -> put in Merkle Tree -> get root hash -> put in block
// Leaves are added in the order that they are hashed (sequential order)
// Note: Merkle Tree is a binary tree, a new node is created with the old root as its left child and the new hash as its right child.
// Step 2:
// Blocks -> get hash -> put in Merkle tree -> get root hash -> put in chain

// Can only add block to chain if it is accepted by the network
