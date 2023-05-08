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

	//blockchain "github.com/Lqvendar/blockchain/blockchain"
	"github.com/Lqvendar/blockchain/blockchain"
	connection "github.com/Lqvendar/blockchain/node"
)

func main() {
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

	chain := blockchain.NewBlockChain()
	currentBlock := blockchain.MakeBlock(chain.GetRoot().GetHash())

	time.Sleep(8 * time.Second)

	fmt.Println("\n--- Welcome to Blockchain! ---\n")

	for {
		fmt.Printf("-----\n\nWhat would you like to do?\n\n1. Add a transaction\n2. View current chain\n\n-----\n\nType option: ")
		reader := bufio.NewScanner(os.Stdin)
		reader.Scan()
		option := reader.Text()

		if option == "2" {
			//invalid memory address
			fmt.Println("\n" + chain.String())

		} else if option == "1" {
			fmt.Printf("\nType transaction data: ")
			reader2 := bufio.NewScanner(os.Stdin)
			reader2.Scan()
			data := reader.Text()

			if data == "" {
				fmt.Println("\n\nData cannot be empty!")
			} else {
				byteData := make([]byte, len(data))

				for i := 0; i < len(data); i++ {
					byteData[i] = data[i]
				}

				transaction := blockchain.Transaction{
					Sender:    []byte("s" + strconv.Itoa(myID)),
					Recipient: []byte("r" + strconv.Itoa(myID)),
					Timestamp: time.Now().UnixNano(),
					Data:      byteData,
				}

				currentBlock.AddTransaction(transaction)
				fmt.Println("\nTransaction added!\n")
			}

		} else {
			fmt.Println("Invalid input! Please select one of the valid options.\n")
		}

		// if block is full
		if len(currentBlock.GetDataList()) == blockchain.GetMax() {
			fmt.Printf("\n\n-----\n\nThe local block has been filled!")

			currentBlock.Mine()
			chain.AddBlock(currentBlock)

			fmt.Printf("Sending block to all nodes...\n\n-----\n\n")
			node.SendBlock(currentBlock)

			// clear current block
			currentBlock = blockchain.MakeBlock(chain.GetRoot().GetHash())

		}
	}
}
