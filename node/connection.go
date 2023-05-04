package node

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"os"
	"sync"

	blockchain "github.com/Lqvendar/blockchain/blockchain"
)

type ServerConnection struct {
	serverID      int
	address       string
	rpcConnection *rpc.Client
}

type Node struct {
	ID        int
	Self      ServerConnection
	peerNodes []ServerConnection

	// local copy of blockchain
	localChain blockchain.BlockChain

	mutex sync.Mutex
	//wg    sync.WaitGroup
}

// Receives blocks and adds them to local blockchain.
// Sets parent hash of block to match current root hash of chain, because the new block will become the new root.
func (node *Node) ReceiveBlock(block *blockchain.Block, reply *string) error {
	fmt.Println("--------------------------------------------")
	fmt.Println("Block received!!")
	fmt.Println("--------------------------------------------")

	block.SetBlockParentHash(node.localChain.GetRoot().GetHash())

	node.localChain.AddBlock(block)

	return nil
}

// Send blocks to peers.
// Mine the block and add it to the chain.
// Calls ReceiveBlock on all peers and prints messages to console (will change to logs)
// Currently: no exported fields error on type block
func (node *Node) SendBlock(block *blockchain.Block, reply *string) error {

	for i, peerNode := range node.peerNodes {
		fmt.Printf("-- Sending block to node %d!\n", i)

		// type blockchain.Block has no exported fields?
		err := peerNode.rpcConnection.Call("node.ReceiveBlock", block, &reply)

		fmt.Println(err)

		if err != nil {
			fmt.Println("Block failed to send!")
		} else {
			fmt.Printf("Block successfully sent to node %d!\n", i)
		}
	}
	return nil
}

func (node *Node) NodeChainToString() string {
	return node.localChain.String()
}

func MakeNode(i int) *Node {
	node := new(Node)
	node.ID = i
	node.Self = ServerConnection{serverID: i}
	node.localChain = *blockchain.NewBlockChain()

	return node
}

func (node *Node) ConnectNodes() error {
	rpc.HandleHTTP()

	node.mutex.Lock()
	selfAddress := node.Self.address
	peerNodes := node.peerNodes
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

			node.peerNodes[index] = ServerConnection{rpcConnection: client}
			fmt.Println("Connected to " + address)

		}(peerNode.address, i)
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
			peerNodes = append(peerNodes, ServerConnection{address: line}) //read only, no need to use mutex
		} else {
			node.Self = ServerConnection{address: line}
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
			peerNodes[i].serverID = i + 1
		} else {
			peerNodes[i].serverID = i
		}
	}

	node.peerNodes = peerNodes
}

/*
arguments := os.Args
    if len(arguments) == 1 {
        fmt.Println("Please provide the server IP and port as host:port.")
        return
    }
func main() {
    fmt.Print("Enter message to be sent: ")
    reader := bufio.NewReader(os.Stdin)

        // read in string from command line
    input, err := reader.ReadString('\n')
    if err != nil {
        fmt.Println("An error occurred while reading input. Please try again.", err)
        return
    }

    // remove the delimeter from the string
    input = strings.TrimSuffix(input, "\n")

        // add to block with input provided
}
*/
