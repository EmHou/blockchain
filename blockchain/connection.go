package blockchain

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"os"
	"sync"

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

	Block      *Block      // block that is being filled up with transactions has NOT been added to the chain yet
	LocalChain *BlockChain // local copy of blockchain

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

type TransactionArg struct {
	Sender    []byte
	Recipient []byte
	Timestamp int64
	Data      []byte

	BlockTimestamp int64 // time that the block was created
}

type TransactionReply struct {
	Success bool
}

func (node *Node) GetSelfAddress() string {
	return node.Self.address
}

// RPC that allows a node to receive a block from another node
// Receives the block data from another node and adds it to its own chain
// If the block is valid, it will add it to its own chain
func (node *Node) ReceiveBlock(args BlockArg, reply *BlockReply) error {
	fmt.Println("--------------------------------------")
	if args.Nonce == 0 && args.DataList == nil {
		// means that empty block was sent 
		addBlock := MakeAddBlock(args.Timestamp, args.Hash, args.Nonce, args.DataList)
		node.Block = addBlock
		fmt.Println(">>> Saving empty block")

	} else {
		// means that the block is full and needs to be added to the chain

		parentBlockHash := node.LocalChain.GetRoot().GetParentBlockHash()
		addBlock := MakeAddBlock(args.Timestamp, parentBlockHash, args.Nonce, args.DataList)
		// create a new Merkle Tree from the list of transactions

		node.wg.Add(1)
		// Nonce should be correct
		go func() {
			defer node.wg.Done()
			err := node.LocalChain.AddConsensusBlock(addBlock, args.Hash)

			if err != nil {
				fmt.Println("RPC >>> Error adding full block to chain")
				reply.Success = false
			} else {
				fmt.Printf("RPC >>> Successfully added full block to chain. Hash: %x\n", addBlock.GetHash())
				node.Block = nil // resets block if added to the chain successfully
				reply.Success = true
			}
		}()
		node.wg.Wait()

		fmt.Println("RPC >>> Successfully added full block to chain")
	}
	fmt.Print("-----\n\nWhat would you like to do?\n\n1. Send a transaction\n2. View hash of local chain\n\n-----\n\nType option: \n")
	return nil
}

// Takes in a block and calls ReceiveBlock on all peer nodes, passing it as an argument.
func (node *Node) SendBlock(block *Block) {
	arg := &BlockArg{
		Nonce:     block.GetNonce(),
		Timestamp: block.GetTimestamp(),
		Hash:      block.GetHash(),
		DataList:  block.GetDataList(),
	}

	for _, peer := range node.peerNodes {
		if peer.rpcConnection != nil {

			node.wg.Add(1)
			go func(peer ServerConnection) {
				defer node.wg.Done()

				var reply BlockReply
				peer.rpcConnection.Call("Node.ReceiveBlock", arg, &reply)

				if reply.Success {
					fmt.Println("Response >>> block was received successfully!")
				} else {
					fmt.Println("Response >>> block was not received successfully!")
				}
			}(peer)
			node.wg.Wait()
		}
	}
}

// RPC that allows a node to receive a transaction from another node
// Receives the transaction data from another node and adds it to its own block
func (node *Node) ReceiveTransaction(args TransactionArg, reply *TransactionReply) error {
	fmt.Println("--------------------------------------")
	newTransaction := &Transaction{
		Sender:    args.Sender,
		Recipient: args.Recipient,
		Timestamp: args.Timestamp,
		Data:      args.Data,
	}

	// Needs to intialise a new block if it doesn't have one
	// with the same timestamp as the se block
	if node.LocalChain.GetBlockListLen() == 1 && node.Block == nil {
		fmt.Println(">>> Creating second block!")
		node.Block = MakeAddBlock(args.BlockTimestamp, node.LocalChain.genesis.GetHash(), 0, nil)
	}

	err := node.Block.Add(*newTransaction)

	if err != nil {
		fmt.Println("RPC >>> Error adding transaction to block")
		reply.Success = false
	} else {
		hash, _ := newTransaction.CalculateHash()
		fmt.Printf("RPC >>> Successfully added transaction to block. Hash: %x\n", hash)
		reply.Success = true
	}
	fmt.Print("--------------------------------------\n\nWhat would you like to do?\n\n1. Send a transaction\n2. View hash of local chain\n\n-----\n\nType option: \n")

	return nil
}

// Sends a transaction to all peer nodes
func (node *Node) SendTransaction(transaction Transaction) {
	arg := &TransactionArg{
		Sender:    transaction.Sender,
		Recipient: transaction.Recipient,
		Timestamp: transaction.Timestamp,
		Data:      transaction.Data,

		BlockTimestamp: node.Block.GetTimestamp(),
	}

	for _, peer := range node.peerNodes {
		if peer.rpcConnection != nil {

			node.wg.Add(1)
			go func(peer ServerConnection) {
				defer node.wg.Done()

				var reply TransactionReply
				peer.rpcConnection.Call("Node.ReceiveTransaction", arg, &reply)

				if reply.Success {
					fmt.Println("Response >>> transaction was received successfully!")
				} else {
					fmt.Println("Response >>> transaction was not received successfully!")
				}
			}(peer)

			node.wg.Wait()
		}
	}
}

// Returns the local chain as a string
func (node *Node) NodeChainToString() string {
	return node.LocalChain.String()
}

// MakeNode creates a new node with the given ID
func MakeNode(i int) *Node {
	node := new(Node)
	node.ID = i
	node.Self = ServerConnection{serverID: i}

	return node
}

// Connects to all the peer nodes listed in the config file
// Allows for asynchronous connecting
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

// Reads the config file and stores the addresses of the peer nodes in the node struct
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
