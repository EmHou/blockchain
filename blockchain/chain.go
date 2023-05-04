package blockchain

import (
	"bytes"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/cbergoon/merkletree"
)

// BlockChain structure links together blocks in a Merkle Tree.
// currentBlock		block that is being filled up with transactions has NOT been added to the chain yet
// chain			chain of blocks that have been added to the chain
type BlockChain struct {
	root      *Block
	genesis   *Block
	chain     *merkletree.MerkleTree
	blockList []merkletree.Content

	wg    sync.WaitGroup
	mutex sync.Mutex
}

// NewBlockChain creates a new blockchain with a genesis block.
func NewBlockChain() *BlockChain {
	genesis := MakeGenesisBlock()
	list := []merkletree.Content{genesis}

	tree, err := merkletree.NewTree(list)
	if err != nil {
		log.Fatal(err)
	}

	blockChain := &BlockChain{
		root:      genesis,
		genesis:   genesis,
		chain:     tree,
		blockList: list,
	}

	return blockChain
}

// Gets the root of the blockchain.
func (blockChain *BlockChain) GetRoot() *Block {
	return blockChain.root
}

// AddBlock adds a block to the blockchain.
// Does not allow block to be added if not full.
// Checks if the hash of the to-be-added block and parent block exists.
// checks if the root hash (of the BlockChain Merkle Tree) is equal to the parent block hash of the to-be-added block.
func (blockChain *BlockChain) AddBlock(block *Block) {
	blockChain.wg.Add(1)
	go func() {
		defer blockChain.wg.Done()

		block.Mine()
	}()
	blockChain.wg.Wait()

	rootHash := blockChain.root.GetHash()
	blockHash := block.GetParentBlockHash()

	if len(block.GetDataList()) < GetMax() {
		log.Println("Block is not full, cannot add to chain.")
		fmt.Println()

		// blockHash exists         parentHash Exists      rootHash == parentBlockHash
	} else if block.GetHash() != nil && rootHash != nil && bytes.Equal(blockHash, rootHash) {
		blockChain.blockList = append(blockChain.blockList, block)
		blockChain.root = block

		blockChain.chain.RebuildTreeWith(blockChain.blockList) // rebuilds chain and sets blockChain.chain to the new chain
	} else {
		log.Println("Block hash does not match root hash, or block hash does not exist.")
		fmt.Println()
	}
}

// Asycnchronously runs the verification of the blockchain every 300 milliseconds.
// This is to ensure no malicious blocks are added to the chain.
func (blockChain *BlockChain) RunVerification() {
	timer := time.NewTimer(300 * time.Millisecond)

	for {
		<-timer.C
		check, err := blockChain.chain.VerifyTree()

		if !check {
			log.Fatal("BlockChain is invalid.")
		}

		if err != nil {
			log.Fatal(err)
		}

		timer.Reset(300 * time.Millisecond)
	}
}

// String representation of the blockchain.
func (blockChain *BlockChain) String() string {
	str := "***BlockChain***\n"

	for _, content := range blockChain.blockList {
		block, _ := content.(*Block)
		str += block.String()
	}

	return str
}
