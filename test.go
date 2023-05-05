package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"time"
	
	connection "github.com/Lqvendar/blockchain/node"

	"sync"

	blockchain "github.com/Lqvendar/blockchain/blockchain"
	"github.com/cbergoon/merkletree"
)

var t = blockchain.Transaction{
	Sender:    []byte("sender"),
	Recipient: []byte("recipient"),
	Timestamp: 1682708458208064000,
	Data:      []byte("data"),
}

var t2 = blockchain.Transaction{
	Sender:    []byte("sender"),
	Recipient: []byte("recipient"),
	Timestamp: 1682708458208064000,
	Data:      []byte("data"),
}

var t3 = blockchain.Transaction{
	Sender:    []byte("sender"),
	Recipient: []byte("recipient"),
	Timestamp: 0,
	Data:      []byte("data"),
}

var trans1 = blockchain.Transaction{
	Sender:    []byte("s1"),
	Recipient: []byte("r1"),
	Timestamp: 1682708458208064001,
	Data:      []byte("data1"),
}
var trans2 = blockchain.Transaction{
	Sender:    []byte("s2"),
	Recipient: []byte("r2"),
	Timestamp: 1682708458208064002,
	Data:      []byte("data2"),
}
var trans3 = blockchain.Transaction{
	Sender:    []byte("s3"),
	Recipient: []byte("r3"),
	Timestamp: 1682708458208064003,
	Data:      []byte("data3"),
}
var trans4 = blockchain.Transaction{
	Sender:    []byte("s4"),
	Recipient: []byte("r4"),
	Timestamp: 1682708458208064004,
	Data:      []byte("data4"),
}

// Tests the getters of a block
func testGetters() {
	fmt.Println("Blockchain")
	fmt.Println("-----------")

	block := blockchain.MakeBlock([]byte{})

	// Testing the getters of a block
	fmt.Println("Testing blockheader (not including parentblockchash and hash)")
	fmt.Println("Nonce: " + strconv.Itoa(int(block.GetNonce())))
	fmt.Println("Timestamp: " + strconv.Itoa(int(block.GetTimestamp())))

	fmt.Println()
}

// Tests Equals() of transactions
func testEquals() {
	fmt.Println("-----------")

	fmt.Println("Testing Equals() of transactions")
	boolean, _ := t.Equals(t2)
	fmt.Println("Should return true: " + strconv.FormatBool(boolean))
	boolean2, _ := t.Equals(t3)
	fmt.Println("Should return false: " + strconv.FormatBool(boolean2))

	fmt.Println()
}

// Tests transaction Merkle Tree
// Creates Merle Tree of transactions
// Verifies hashes of entire Merkle Tree
// Verifies a specific transaction is in the Merkle Tree
// Prints out string representation of the Merkle Tree
func testTransactionTree() {
	fmt.Println("-----------")

	// Testing how to append transactions to data (MerkleTree)
	var list []merkletree.Content
	list = append(list, t)
	list = append(list, t3)
	list = append(list, t2)

	// create a new Merkle Tree from the list of transactions
	tree, err := merkletree.NewTree(list)
	if err != nil {
		log.Fatal(err)
	}

	// Verify the entire tree (hashes for each node) is valid
	vt, err := tree.VerifyTree()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Testing VerifyTree()")
	fmt.Println("Should return true: ", vt)
	fmt.Println()

	//Verify a specific content in in the tree
	vc, err := tree.VerifyContent(t)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("-----------")
	fmt.Println("Testing VerifyContent()...")
	fmt.Println("Should return true (because \"t\" should be in the MerkleTree): ", vc)
	fmt.Println()

	// Prints out in order: leaf (bool), duplicate (bool), hash ([]byte), content (interface)
	// Prints out
	fmt.Println("-----------")
	fmt.Println("Printing out string representation of the MerkleTree of transactions")
	fmt.Println(tree.String())
	fmt.Println("Getting MerkleRoot of the MerkleTree")
	fmt.Println("Should return the hash : ", tree.MerkleRoot())

	fmt.Println()
}

// Tests GetMerklePath() and RebuildTreeWith()
// This is to understand how Merkle trees work and how to use the functions
func testMerkleTree() {
	fmt.Println("-----------")
	var newList []merkletree.Content
	newList = append(newList, trans1)
	newList = append(newList, trans2)
	newList = append(newList, trans3)

	// create a new Merkle Tree from the list of transactions
	newTree, errrrrrr := merkletree.NewTree(newList)
	if errrrrrr != nil {
		log.Fatal(errrrrrr)
	}

	// testing verify
	verify, err := newTree.VerifyTree()
	if err != nil {
		log.Fatal(err)
	}

	// Testing GetMerklePath()
	fmt.Println("Verifying newTree: " + strconv.FormatBool(verify))
	fmt.Println("String representation of newTree: " + newTree.String())
	fmt.Println()

	fmt.Println("Testing GetMerklePath() of trans1")
	path, indices, _ := newTree.GetMerklePath(trans1)
	fmt.Println("Should return the path of trans1: ", path, indices)

	fmt.Println("Testing GetMerklePath() of trans2")
	path2, indices2, _ := newTree.GetMerklePath(trans2)
	fmt.Println("Should return the path of trans2: ", path2, indices2)

	fmt.Println("Testing GetMerklePath() of trans3")
	path3, indices3, _ := newTree.GetMerklePath(trans3)
	fmt.Println("Should return the path of trans3: ", path3, indices3)
	fmt.Println()

	// Testing RebuildTreeWith() with trans4
	// RebuildTreeWith() just rebuilds the tree after adding a new transaction and will keep the same hashes for the previous leaves
	fmt.Println("-----------")
	newList = append(newList, trans4)
	// create a new Merkle Tree from the list of transactions
	evenNewerTree, err100 := merkletree.NewTree(newList)
	if err100 != nil {
		log.Fatal(err100)
	}
	fmt.Println("String representation of evenNewerTree: ")
	fmt.Println(evenNewerTree.String())

	fmt.Println()
}

// Tests adding transactions to a block
// Tests if maximum transactions in a block is 1 (because of SetMax(1))
// Testing Mine() (finding nonce and setting hash)
func testAddTransactionsAndMine() {
	fmt.Println("-----------")

	newBlock := blockchain.MakeBlock([]byte{0})
	blockchain.SetMax(1)
	newBlock.AddTransaction(trans1)
	newBlock.AddTransaction(trans2) // should not be added
	fmt.Println("Testing AddTransaction() to a block")
	fmt.Println("String representation of newBlock: ")
	fmt.Println(newBlock.GetData().String()) // will return a duplicate of trans 1

	fmt.Println()
	fmt.Println("-----------")

	fmt.Println("Testing Mine()")
	nonce, _ := newBlock.Mine()
	fmt.Println("Nonce: " + strconv.Itoa(nonce))
	fmt.Printf("Set hash in block: %x", newBlock.GetHash())

	fmt.Println()
	fmt.Println()
}

func testNewBlockChain() {
	fmt.Println("-----------")
	fmt.Println("Testing NewBlockChain()")
	chain := blockchain.NewBlockChain()
	fmt.Println("String representation of chain: ")
	fmt.Println(chain.String())

	fmt.Println("-----------")
	fmt.Println("Testing AddBlock()")
	blockchain.SetMax(7)
	block := blockchain.MakeBlock([]byte{0})
	block.AddTransaction(trans1)
	block.AddTransaction(trans2)
	block.AddTransaction(trans3)
	block.AddTransaction(trans4)
	block.AddTransaction(t)
	block.AddTransaction(t2)
	//block.AddTransaction(t3)

	chain.AddBlock(block)
	fmt.Println()
	// also shouldn't add because doesn't have parent hash
	fmt.Println("Should not add to chain because not at max transactions:")
	fmt.Println(chain.String())
}

func testAddingCorrectBlock() {
	fmt.Println("-----------")
	chain := blockchain.NewBlockChain()
	fmt.Println("String representation of chain (with genesis block): ")
	fmt.Println(chain.String())
	fmt.Println()
	fmt.Println("Testing AddBlock() with correct block")
	goodBlock := blockchain.MakeBlock(chain.GetRoot().GetHash())
	fmt.Print("Parent hash: ")
	fmt.Println(chain.GetRoot().GetHash())
	goodBlock.AddTransaction(trans1)
	goodBlock.AddTransaction(trans2)
	goodBlock.AddTransaction(trans3)
	goodBlock.AddTransaction(trans4)
	goodBlock.AddTransaction(t)
	goodBlock.AddTransaction(t2)
	goodBlock.AddTransaction(t3)

	chain.AddBlock(goodBlock)

	fmt.Println("Should add to chain because correct block:")
	fmt.Println(chain.String())
	fmt.Println()
}

func testMining(difficulty int) {
	testMineBlock := blockchain.MakeBlock([]byte{0})
	testMineBlock.AddTransaction(trans1)

	testMineBlock.TestPOW(difficulty)

	testMineBlock.TestPrintMine()
}

func testConnection() {
	arguments := os.Args

	myID, err := strconv.Atoi(arguments[1])
	node := connection.MakeNode(myID)
	if err != nil {
		log.Fatal(err)
	}

	err = rpc.Register(node)
	fmt.Println("Node " + strconv.Itoa(myID) + " up!")
	if err != nil {
		log.Fatal("error registering the RPCs\n", err)
	}

	node.ReadClusterConfig("nodes.txt")
	// nodes connect now
	node.ConnectNodes()
	time.Sleep(8 * time.Second)

	// testing creation of a chain
	chain := blockchain.NewBlockChain()

	// testing creation of a block
	block := blockchain.MakeBlock(chain.GetRoot().GetHash())
	block.AddTransaction(trans1)
	block.AddTransaction(trans2)
	block.AddTransaction(trans3)
	block.AddTransaction(trans4)
	block.AddTransaction(t)
	block.AddTransaction(t2)
	//block.AddTransaction(t3)


	chain.AddBlock(block)
	node.SendBlock(block)
}

/*
// This is for the Blockchain presenation
func main() {
	// Testing the creation of a block
	// Nonce set to 0
	// Hash is not set yet because it will be set when mining
	//testGetters()

	// Added three transactions to the block
	// VerifyTree() verifies the hashes of the Merkle Tree
	// VerifyContent() verifies if a specific transaction is in the Merkle Tree
	// String representation of Merkle Tree
	// first boolean value: if the content stored is a leaf
	// second boolean value: if it is a duplicate
	// [] of numbers is the hash
	// info of the block header
	//t()

	/t()

	// Setting max number of transactions in a block to 1 just for demonstration purposes
	// Attempt to add two transactions into the block
	// Only one transaction will be added because the max is 1
	// In the string representation of the block, it will show the duplicate transaction (second boolean is true)
	// Attempts to mine the block. After it's done mining, it will set the block's hash and nonce.
	//testAddTransactionsAndMine()

	// Making a block and adding a transaction (max is still 1)
	// Prints out all the hashes that it tests with a changing nonce (nonce increases by 1 each time)
	//testMining(18)
}
*/

/*
func main() {
	// Testing the getters of a block
	testGetters()

	// Testing Equals() of transactions
	testEquals()

	// Testing transaction Merkle Tree
	// Creates Merle Tree of transactions
	// Verifies hashes of entire Merkle Tree
	// Verifies a specific transaction is in the Merkle Tree
	// Prints out string representation of the Merkle Tree
	t()

	// Testing GetMerklePath() and RebuildTreeWith()
	// This is to understand how Merkle trees work and how to use the functions
	testMerkleTree()

	// Testing AddTransaction() to a block and Run()
	testAddTransactionsAndMine()
	chain := testNewBlockChain()

t(&chain) // & refers to where the variable is located in memory

	var wg sync.WaitGroup
	wg.Add(1)
	go chain.RunVerification()
	wg.Wait()
}
*/

func main() {
	testConnection()

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
