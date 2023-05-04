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
	connection "github.com/Lqvendar/blockchain/node"
)

func main() {

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

	// for test-sending blocks, will remove //

	// var trans1 = blockchain.Transaction{
	// 	Sender:    []byte("s1"),
	// 	Recipient: []byte("r1"),
	// 	Timestamp: 1682708458208064001,
	// 	Data:      []byte("data1"),
	// }

	// var trans2 = blockchain.Transaction{
	// 	Sender:    []byte("s2"),
	// 	Recipient: []byte("r2"),
	// 	Timestamp: 1682708458208064002,
	// 	Data:      []byte("data2"),
	// }

	// blockchain.SetMax(1)

	// newBlock := blockchain.MakeBlock([]byte{0})
	// newBlock.AddTransaction(trans1)
	// newBlock.Mine()

	//newBlock.GetHash()

	// newBlock2 := blockchain.MakeBlock([]byte{0})
	// newBlock2.AddTransaction(trans2)
	// newBlock2.Mine()

	//newBlock2.GetHash()

	reply := ""

	//node.SendBlock(newBlock, &reply)
	//node.SendBlock(newBlock2, &reply)

	// Test-sending ends here //

	/*
		Basic order of events:
		1. Ask for user input and store in var DATA
		2. Create transaction with this data and proper info
		3. Add transaction to block... continue until max reached
		4. Send block to all nodes
		5. Nodes receive blocks, add to their own blockchains
	*/

	// create new scanner to read input from command line
	// transaction data goes here
	for {
		fmt.Print("Type 1 to send data, type 2 to view current chain: ")
		reader := bufio.NewScanner(os.Stdin)
		reader.Scan()
		option := reader.Text()

		if option == "2" {
			fmt.Print(node.NodeChainToString())
		} else if option == "1" {
			fmt.Println("What would you like to send?: ")
			reader2 := bufio.NewScanner(os.Stdin)
			reader2.Scan()
			data := reader.Text()
			byteData := make([]byte, len(data))

			for i := 0; i < len(data); i++ {
				byteData[i] = data[i]
			}

			// make current block and add transaction
			currentBlock := blockchain.MakeBlock([]byte{0})

			// this is for user input:

			// make transaction from data

			// transaction := blockchain.Transaction{
			// 	Sender:    []byte("s1"),
			// 	Recipient: []byte("r1"),
			// 	Timestamp: time.Now().UnixNano(),
			// 	Data:      byteData,
			// }
			// currentBlock.AddTransaction(transaction)

			// user input ends here

			// for non-user input testing:
			currentBlock.AddTransaction(t)
			currentBlock.AddTransaction(t2)
			currentBlock.AddTransaction(t3)
			currentBlock.AddTransaction(trans1)
			currentBlock.AddTransaction(trans2)
			currentBlock.AddTransaction(trans3)
			currentBlock.AddTransaction(trans4)

			// maybe add this to SendBlock RPC instead
			currentBlock.Mine()

			// test get peer nodes

			// fmt.Println("--------------------------------------------")
			// fmt.Println("PEER NODES:")
			// fmt.Print(node.GetPeerNodes())
			// fmt.Println("--------------------------------------------")

			// test SendBlock
			fmt.Println("--------------------------------------------")
			fmt.Println("BLOCK SENT")
			fmt.Println("--------------------------------------------")
			node.SendBlock(currentBlock, &reply)
			fmt.Println("--------------------------------------------")

			// fmt.Println("SECOND BLOCK SENT")
			// fmt.Println("--------------------------------------------")
			// node.SendBlock(currentBlock, &reply)
			// fmt.Println("--------------------------------------------")

			//test ReceiveBlock
			// fmt.Println("--------------------------------------------")
			// fmt.Println("BLOCK RECEIVED")
			// fmt.Println("--------------------------------------------")
			// node.ReceiveBlock(currentBlock, &reply)
			// fmt.Println("--------------------------------------------")

			// fmt.Println("SECOND BLOCK RECEIVED")
			// fmt.Println("--------------------------------------------")
			// node.ReceiveBlock(newBlock2, &reply)
			// fmt.Println("--------------------------------------------")

		} else {
			fmt.Println("Invalid entry!")
		}
	}
}
