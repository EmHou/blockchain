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
	// wg    sync.WaitGroup
}

// ReceiveBlock: The block sent to/received by the node.
type ReceiveBlockArg struct {
	ReceiveBlock *blockchain.Block
}

// Success: Marks whether or not the block was successfully added to the chain.
type ReceiveBlockReply struct {
	Success bool
}

// ReceiveBlock RPC. Takes in ReceiveBlockArg (block received), sets parent hash, and adds to blockchain.
// Replies success if no errors occurred.
func (node *Node) ReceiveBlock(arg ReceiveBlockArg, reply *ReceiveBlockReply) error {
	arg.ReceiveBlock.SetBlockParentHash(node.localChain.GetRoot().GetHash())
	node.localChain.AddBlock(arg.ReceiveBlock)
	reply.Success = true
	return nil
}

// Takes in a block and calls ReceiveBlock on all peer nodes, passing it as an argument.
func (node *Node) SendBlock(block *blockchain.Block) {
	arg := new(ReceiveBlockArg)
	arg.ReceiveBlock = block

	for i := range node.peerNodes {
		result := new(ReceiveBlockReply)

		// this is fine
		// fmt.Printf("============STACK=============")
		// debug.PrintStack()
		// fmt.Printf("==============================")

		serverCall := node.peerNodes[i].rpcConnection.Go("Node.ReceiveBlock", arg, &result, nil)

		// doesn't reach this
		// debug.PrintStack()

		<-serverCall.Done

		if result.Success {
			fmt.Println("Block sent successfully!")
		}
	}
}

// old //

// Receives blocks and adds them to local blockchain.
// Sets parent hash of block to match current root hash of chain, because the new block will become the new root.
// func (node *Node) ReceiveBlock(block *blockchain.Block, reply *string) error {
// 	fmt.Println("--------------------------------------------")
// 	fmt.Println("Block received!!")
// 	fmt.Println("--------------------------------------------")

// 	block.SetBlockParentHash(node.localChain.GetRoot().GetHash())

// 	node.localChain.AddBlock(block)

// 	return nil
// }

// Send blocks to peers.
// Mine the block and add it to the chain.
// Calls ReceiveBlock on all peers and prints messages to console (will change to logs)
// stack overflow error
// func (node *Node) SendBlock(block *blockchain.Block) {
// 	var wg sync.WaitGroup
// 	var i = 0

// 	arg := new(ReceiveBlockArg)
// 	arg.ReceiveBlock = block

// 	for _, peerNode := range node.peerNodes {
// 		wg.Add(1)

// 		go func(peerNode ServerConnection) {
// 			defer wg.Done()

// 			result := new(ReceiveBlockReply)

// 			err := peerNode.rpcConnection.Go("Node.ReceiveBlock", arg, &result, nil)

// 			if err != nil {
// 				fmt.Println("Block failed to send!")
// 			} else {
// 				fmt.Printf("Block successfully sent to node %d!\n", i)
// 			}
// 			i++

// 		}(peerNode)
// 	}

// 	wg.Wait()

// 	fmt.Println("-- Sent block to nodes!")
// }

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
