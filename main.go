// main function

package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"time"

	blockchain "github.com/Lqvendar/blockchain/blockchain"
	//"github.com/Lqvendar/blockchain/blockchain"
	// connection "github.com/Lqvendar/blockchain/node"
)

func main() {
	arguments := os.Args

	myID, err := strconv.Atoi(arguments[1])
	node := blockchain.MakeNode(myID)
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

	chain := blockchain.NewBlockChain()

	time.Sleep(8 * time.Second)

	fmt.Println("\n--- Welcome to Blockchain! ---\n")

	for {
		fmt.Printf("-----\n\nWhat would you like to do?\n\n1. Send a transaction\n\n2. View current block\n\n-----\n\nType option: ")
		reader := bufio.NewScanner(os.Stdin)
		reader.Scan()
		option := reader.Text()

		if option == "1" {
			fmt.Println(">>> Type transaction data: ")
			reader2 := bufio.NewScanner(os.Stdin)
			reader2.Scan()
			data := reader.Text()

			if data == "" {
				fmt.Println(">>> Data cannot be empty!")
			} else {
				fmt.Println(">>> Creating transaction!")
				transaction := blockchain.MakeTransaction(node.GetSelfAddress(), "recipient", time.Now().UnixNano(), data)

				// Check if block list length is 1 (empty besides genesis) and create block if so
				if chain.GetBlockListLen() == 1 {
					fmt.Println(">>> Creating first block!")
					node.Block = blockchain.MakeBlock(chain.GetRoot().GetHash())
					
					node.Block.AddTransaction(*transaction, chain, node)

				} else { // Else, add and send transaction
					node.Block.AddTransaction(*transaction, chain, node)
					fmt.Println(">>> Added transaction to local chain!")
					node.SendTransaction(*transaction)
					fmt.Println(">>> Sent transaction to all nodes!")
				}
			}
		} else {
			fmt.Println(">>> Invalid input! Please select one of the valid options.\n")
		}

	}
}

// byteData := make([]byte, len(data))

// for i := 0; i < len(data); i++ {
// 	byteData[i] = data[i]
// }

// transaction := blockchain.Transaction{
// 	Sender:    []byte("s" + strconv.Itoa(myID)),
// 	Recipient: []byte("r" + strconv.Itoa(myID)),
// 	Timestamp: time.Now().UnixNano(),
// 	Data:      byteData,
// }
