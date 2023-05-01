package node

import (
	"log"
	"net/http"
)

func (node *RaftNode) connectNodes() error {
	rpc.HandleHTTP()

	node.mutex.Lock()
	selfAddress := node.self.Address
	peerNodes := node.peerNodes
	node.mutex.Unlock()

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

		}(peerNode.Address, i)
	}
	return nil
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
