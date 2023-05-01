package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"sync"
	"time"
	//blockchain "github.com/Lqvendar/blockchain/blockchain"
)

type ServerConnection struct {
	ServerID      int
	Address       string
	RpcConnection *rpc.Client
}

type Node struct {
	ID        int
	Self      ServerConnection
	PeerNodes []ServerConnection

	mutex sync.Mutex
	//wg    sync.WaitGroup
}

// needs a function to register to RPC?
func (node *Node) Test(arguments Node, reply *Node) error {
	fmt.Println("test")
	return nil
}

func MakeNode(i int) *Node {
	node := new(Node)
	node.ID = i
	node.Self = ServerConnection{ServerID: i}

	return node
}

func (node *Node) ConnectNodes() error {
	rpc.HandleHTTP()

	node.mutex.Lock()
	selfAddress := node.Self.Address
	peerNodes := node.PeerNodes
	node.mutex.Unlock()

	fmt.Println(selfAddress)

	go http.ListenAndServe(selfAddress, nil)
	log.Printf("Serving rpc on: " + selfAddress)

	for i, peerNode := range peerNodes {
		//node.wg.Add(1)
		go func(address string, index int) {
			//defer node.wg.Done()
			client, err := rpc.DialHTTP("tcp", address)

			for err != nil {
				client, err = rpc.DialHTTP("tcp", address)
			}

			node.PeerNodes[index] = ServerConnection{RpcConnection: client}
			fmt.Println("Connected to " + address)

		}(peerNode.Address, i)
	}
	return nil
}

func (node *Node) ReadClusterConfig(filename string) {
	nodeIndex := node.ID

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file: " + filename)
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var peerNodes []ServerConnection

	index := 0
	for scanner.Scan() {
		line := scanner.Text()

		if index != nodeIndex { // 0:4040 1:4041 2:4042 3:4043 4:4044
			peerNodes = append(peerNodes, ServerConnection{Address: line}) //read only, no need to use mutex
		} else {
			node.Self = ServerConnection{Address: line}
		}

		index++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: " + filename)
		log.Fatal(err)
	}

	//looping through all the nodes as listed in the config file
	for i := 0; i < len(peerNodes); i++ {
		// creating the IDs for each node, without including its own ID
		if i >= nodeIndex {
			peerNodes[i].ServerID = i + 1
		} else {
			peerNodes[i].ServerID = i
		}
	}

	node.PeerNodes = peerNodes
}

func main() {

	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide cluster information.")
		return
	}

	myID, err := strconv.Atoi(arguments[1])
	api := MakeNode(myID)
	if err != nil {
		log.Fatal(err)
	}

	err = rpc.Register(api)
	fmt.Println("Node " + strconv.Itoa(myID) + " up!")
	if err != nil {
		log.Fatal("error registering the RPCs\n", err)
	}

	api.ReadClusterConfig("nodes.txt")

	// nodes connect now
	api.ConnectNodes()

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
