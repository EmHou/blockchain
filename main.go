// main function

package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"sync"
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

	node.LocalChain = blockchain.NewBlockChain()

	time.Sleep(3 * time.Second)

	fmt.Println("\n--- Welcome to Blockchain! ---\n")

	continueLoop := true

	for continueLoop {
		var wg sync.WaitGroup
		
		fmt.Print("--------------------------------------\n\nWhat would you like to do?\n\n1. Send a transaction\n2. View hash of local chain\n\n-----\n\nType option: \n")
		reader := bufio.NewScanner(os.Stdin)
		reader.Scan()
		option := reader.Text()

		if option == "1" {
			fmt.Println(">>> Enter recipient:")
			reader.Scan()
			recipient := reader.Text()

			fmt.Println(">>> Type transaction data: ")
			reader.Scan()
			data := reader.Text()

			if data == "" {
				fmt.Println(">>> Data cannot be empty!")
				fmt.Println()
			} else {
				wg.Add(1)
				go func() {
					defer wg.Done()
					fmt.Println(">>> Creating transaction!")
					fmt.Println()
					transaction := blockchain.MakeTransaction(node.GetSelfAddress(), recipient, time.Now().UnixNano(), data)

					// Check if block list length is 1 (empty besides genesis) and create block if so
					if node.LocalChain.GetBlockListLen() == 1 && node.Block == nil {
						fmt.Println(">>> Creating main block!")
						fmt.Println()
						node.Block = blockchain.MakeBlock(node.LocalChain.GetRoot().GetHash())
					}

					node.Block.AddTransaction(*transaction, node.LocalChain, node)
					hash, _ := transaction.CalculateHash()
					fmt.Printf(">>> Added transaction to local chain! Hash: %x\n", hash)
					fmt.Println(">>> Sent transaction to all nodes!")
					fmt.Println()
				}()

				wg.Wait()
			}
		} else if option == "2" {
			fmt.Println("Current chain hash: ", node.LocalChain.GetRoot().GetHash())
		} else {
			fmt.Println(">>> Invalid input! Please select one of the valid options.")
			fmt.Println()
		}

		fmt.Println("--------------------------------------")
		fmt.Println("Would you like to continue? (y/n): ")
		reader.Scan()
		option = reader.Text()

		if option == "n" {
			fmt.Println("Exiting...")
			continueLoop = false
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
