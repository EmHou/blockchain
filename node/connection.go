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
	"github.com/cbergoon/merkletree"
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
	localChain *blockchain.BlockChain

	mutex sync.Mutex
	wg    sync.WaitGroup
}

// From its own chain
// Already have from header: parentBlockHash,
// Already have from block: transaction hash of parent block (data)
// Needs: nonce, timestamp, and hash to verify hash
type BlockArg struct {
	Nonce     uint64
	Timestamp int64
	Hash      []byte

	DataList []merkletree.Content
}

type BlockReply struct {
	Success bool
}

func (node *Node) ReceiveBlock(args BlockArg, reply *BlockReply) error {
	
	// Needs to intialise a new blockchain if it doesn't have one
	// with the same timestamp as the genesis block
	if node.localChain == nil{
		node.localChain = blockchain.NewBlockChain()
	}
	
	parentBlockHash := node.localChain.GetRoot().GetParentBlockHash()
	addBlock := blockchain.MakeAddBlock(args.Timestamp, parentBlockHash, args.Nonce, args.DataList)
	// create a new Merkle Tree from the list of transactions

	node.wg.Add(1)
	// Nonce should be correct
	go func() {
		defer node.wg.Done()
		err := node.localChain.AddBlock(addBlock)

		if err != nil {
			reply.Success = false
		} else {
			reply.Success = true
		}
	}()
	node.wg.Wait()
	
	return nil
}

// Takes in a block and calls ReceiveBlock on all peer nodes, passing it as an argument.
func (node *Node) SendBlock(block *blockchain.Block) {
	arg := &BlockArg{
		Nonce:     block.GetNonce(),
		Timestamp: block.GetTimestamp(),
		Hash:      block.GetHash(),
		DataList:  block.GetDataList(),
	}

	for _, peer := range node.peerNodes {
		if peer.rpcConnection != nil {

			go func(peer ServerConnection) {
				var reply BlockReply

				peer.rpcConnection.Call("Node.ReceiveBlock", arg, &reply)

				if reply.Success {
					fmt.Println("Response back from peer node: block was received successfully!")
				}
			}(peer)
		}
	}
}

func (node *Node) NodeChainToString() string {
	return node.localChain.String()
}

func MakeNode(i int) *Node {
	node := new(Node)
	node.ID = i
	node.Self = ServerConnection{serverID: i}

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
