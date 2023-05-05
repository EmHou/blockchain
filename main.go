// main function

package main

import (
	/*
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"time"
	*/

	//blockchain "github.com/Lqvendar/blockchain/blockchain"
	//connection "github.com/Lqvendar/blockchain/node"
)

/*
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

	time.Sleep(8 * time.Second)

	for {
		// fmt.Print("Type 1 to send data, type 2 to view current chain: ")
		// reader := bufio.NewScanner(os.Stdin)
		// reader.Scan()
		// option := reader.Text()

		// if option == "2" {
		// 	fmt.Print(node.NodeChainToString())
		// } else if option == "1" {
		// 	fmt.Println("What would you like to send?: ")
		// 	reader2 := bufio.NewScanner(os.Stdin)
		// 	reader2.Scan()
		// 	data := reader.Text()
		// 	byteData := make([]byte, len(data))

		// 	for i := 0; i < len(data); i++ {
		// 		byteData[i] = data[i]
		// 	}

		// transaction := blockchain.Transaction{
		// 	Sender:    []byte("s1"),
		// 	Recipient: []byte("r1"),
		// 	Timestamp: time.Now().UnixNano(),
		// 	Data:      byteData,
		// }
		// currentBlock.AddTransaction(transaction)

		// user input ends here
	}
}
*/