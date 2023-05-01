package main

/*
import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"time"
	//blockchain "github.com/Lqvendar/blockchain/blockchain"

	connection "github.com/Lqvendar/blockchain/node"
)
*/

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

	// create new scanner to read input from command line
	// transaction data goes here
	for {
		fmt.Println("What would you like to send?: ")
		reader := bufio.NewScanner(os.Stdin)
		reader.Scan()
		data := reader.Text()
		fmt.Println(data)

		// make transaction from data
		// maybe put this somewhere else?
		// transaction := blockchain.Transaction{
		// 	Sender:    api.Self.Address,
		// 	Recipient: api.PeerNodes[1],
		// 	Timestamp: time.Now().UnixNano(),
		// 	Data:      data,
		// }
	}
}
*/
