package main

import (
	"fmt"
	"strconv"

	blockchain "github.com/Lqvendar/blockchain/blockchain"
)

func main() {
	fmt.Println("Blockchain")
	fmt.Println("-----------")
	block := blockchain.MakeBlock([]byte{})

	// Testing the getters of a block
	fmt.Println("Data:")
	fmt.Println("Nonce: " + strconv.Itoa(int(block.GetNonce())))
	fmt.Println("Timestamp: " + strconv.Itoa(int(block.GetTimestamp())))



	// Nodes are connected are connected in the network
	// Genesis block is created
	// Genesis block is added to the empty initialised chain (initialised by Genesis machine)

	// Transaction is called
		// Transaction is added to the block (actually is added in the Merkle Tree)
	// Block fills up and start mining
		// 

}

