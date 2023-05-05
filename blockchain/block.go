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

// Includes helper methods that allow easy access to find ParentBlockHash, Hash, and Nonce.
// Arbitrary number of transactions - in this implemenation, we choose 7 transactions per block.
// Header	contains metaDataof the block
// data		transactions (Merkle Tree)
// pow		Proof of Work algorithm to validate blocks and calculate Nonce to add block to chain
// DataList	transactions that are to added to the block
type Block struct {
	Header   BlockHeader
	Data     *merkletree.MerkleTree
	Pow      *ProofOfWork
	DataList []merkletree.Content
}

// BlockHeader contains metaDataof the block.
// Timestamp 		time when block is added to the chain
// ParentBlockHash	Hash of the previous block
// Hash	 			takes Nonce, Timestamps, ParentBlockHash, and root Hash of the transaction Merkle Tree
// Nonce	 		rand int that is initialised to 0
type BlockHeader struct {
	Timestamp       int64
	ParentBlockHash []byte
	Hash            []byte
	Nonce           uint64
}

func (block *Block) GetTimestamp() int64 {
	return block.Header.Timestamp
}

func (block *Block) GetHeader() BlockHeader {
	return block.Header
}

func (block *Block) GetParentBlockHash() []byte {
	return block.Header.ParentBlockHash
}

func (block *Block) GetHash() []byte {
	return block.Header.Hash
}

func (block *Block) SetHash(Hash []byte) {
	block.Header.Hash = Hash
}

func (block *Block) GetNonce() uint64 {
	return block.Header.Nonce
}

func (block *Block) SetNonce(Nonce uint64) {
	block.Header.Nonce = Nonce
}

func (block *Block) GetData() *merkletree.MerkleTree {
	return block.Data
}

func (block *Block) GetTarget() *big.Int {
	return block.Pow.Target
}

func (block *Block) GetDataList() []merkletree.Content {
	return block.DataList
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
	Hash := sha256.New()
	Data := Block.BlockDataToBytes()

	if _, err := Hash.Write(Data); err != nil {
		return nil, err
	}

	return Hash.Sum(nil), nil
}

// need to set Hash of the block after it is added to the chain
// can only use this method after the block is added to the chain
// because the Data(MerkleTree of transacitons) must be set before the Hash can be set
// Part of the Content interface in MerkleTree Package
// TODO
func (block Block) Equals(other merkletree.Content) (bool, error) {
	return block.Data == other.(Block).Data, nil
}

// NewBlockChain creates a new block and returns the pointer to it.
// Initialises:
// Timestamp		to current time in nanoseconds
// ParentBlockHash	to the Hash of the previous block (gotten from parameter)
// Hash				to empty byte array, to be set after a Nonce is found (after mining when adding to chain)
// Nonce 			to 0
// DataList 		to empty array of type (merkletree.Content)
// Pow				to a new ProofOfWork struct
func MakeBlock(pBlockHash []byte) *Block {
	Header := &BlockHeader{
		Timestamp:       time.Now().UnixNano(),
		ParentBlockHash: pBlockHash,
		Hash:            []byte{},
		Nonce:           0,
	}

	block := &Block{
		Header:   *Header,
		DataList: []merkletree.Content{},
		Pow:      NewPOW(),
	}

	return block
}

// Takes MakeBlock() and passes empty byte array to it
// Has no parent block, so ParentBlockHash is empty (empty byte array)
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
	if len(block.DataList) >= max {
		fmt.Println("Block is full, cannot add more transactions")

	} else {
		block.DataList = append(block.DataList, transaction)

		if block.Data == nil {
			tree, err := merkletree.NewTree(block.DataList)
			if err != nil {
				log.Fatal(err)
			}

			block.Data = tree
		} else {
			err := block.Data.RebuildTreeWith(block.DataList)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

}

// Function used to set parent Hash
// Used when first block is added to the blockchain after the genesis block.
// Otherwise, the parent Hash will be empty for this block, and it will not be added.
func (block *Block) SetBlockParentHash(Hash []byte) {
	block.Header.ParentBlockHash = Hash
}

func (block Block) String() string {
	str := "**Block**\n"
	str += block.Header.String()
	str += "Data(String representation of Transactions Merkle Tree):\n"
	str += block.Data.String() + "\n"

	return str
}

func (Header *BlockHeader) String() string {
	str := "Timestamp: " + strconv.FormatInt(Header.Timestamp, 10) + "\n"
	str += "Parent Block Hash: " + hex.EncodeToString(Header.ParentBlockHash) + "\n"
	str += "Hash: " + hex.EncodeToString(Header.Hash) + "\n"

	return str
}

//-- Notes --//

// Step 1:
// Transaction -> get Hash -> put in Merkle Tree -> get root Hash -> put in block
// Leaves are added in the order that they are Hashed (sequential order)
// Note: Merkle Tree is a binary tree, a new node is created with the old root as its left child and the new Hash as its right child.
// Step 2:
// Blocks -> get Hash -> put in Merkle tree -> get root Hash -> put in chain

// Can only add block to chain if it is accepted by the network
