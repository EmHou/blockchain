# Blockchain: A Matter of Distributed Consensus


## Abstract
Blockchain is a distributed ledger system that has gained significant attention in recent years due to its widespread range of potential applications in multiple industries and systems. In particular, blockchain has found use in cryptocurrency systems such as Bitcoin and Ethereum, which are known for being pseudonymous, immutable, and tamper-proof. This is made possible by blockchain technology, which utilizes structures called blocks to securely store, send, and receive information about data transactions in a distributed network. In this paper, we explore the functionality of blockchain and detail our implementation in Golang using RPCs to send data between nodes. We simulate a blockchain network with three nodes, from which data is entered in a command line interface, added to a block, validated, and sent to other connected nodes. Data is verified using hashes— if a node is to receive a transaction, the command line interface allows us to view the hash of the received record and check if the hash is the same across all nodes. If it is, the process has completed without error, and our implementation succeeds in doing so.  
	
## 1&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Introduction
Blockchain is a distributed system unveiled by Satoshi Nakamoto in his Bitcoin white paper released in 2009 [3]. The technology has gained widespread notoriety for serving as the foundation of cryptocurrencies such as Bitcoin and Ethereum, mainly due to its pseudonymity and security. The blockchain consists of multiple sets of data, known as “blocks'', which are linked together sequentially in such a manner that they are almost impossible to change. In addition, the data within the blocks themselves is heavily encrypted, which makes outside access to the information nearly impossible [1].

The blockchain is entirely decentralized, meaning that it is shared across a very large network of nodes, or computers, and any node within that network is able to execute a transaction, or send and receive data, at any time [5]. Blockchain runs on a peer-to-peer network structure, so a node that wishes to send a transaction to another node will establish a direct connection with it in order to send information [5]. No central body oversees this process, which distinguishes blockchain’s architecture from a client-server model and contributes to its high level of security. 

Nakamoto’s main goal in his proposal and implementation of the blockchain was to establish a distributed system that relied on a “matter of trust”, or more specifically, to ensure that information within said system was tamper-proof [3]. Blockchain’s method of hashing and storing data ensures that nodes are unable to change information within the chain itself without clear detection through validation mechanisms such as verifiable hashes and digital signatures [3].

## 2&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Blockchain
Four major components make up blockchain technology: records, blocks, chains, and consensus [1] (detailed in section 4).

Transactions sent from node to node are stored in structures called records [1], which keep track of transaction information including the timestamp of the transaction, the sender and receiver addresses, digital signatures, and the actual data being sent. Timestamps are used to verify the order in which transactions are sent and received. Doing so allows the blockchain to keep track of transaction history and also contributes to immutability as verification algorithms are able to validate the order in which transactions are added to the chain [5]. Sender and receiver addresses represent the nodes sending and receiving data, but they have other uses as well— for example, in Bitcoin, these would represent the wallets that transactions are sent to and from. Digital signatures verify that sent transactions have not been altered in any way and that the sender intended to send the transaction. When a transaction is sent, it is first hashed to ensure that it cannot be viewed by any outside parties [3]. It is then validated by nodes in the network before, finally, it is added to a structure known as a block [1].

Blocks serve as a mechanism for storing records of transactions that have taken place on the blockchain network. Each block can hold a set number of records, and this maximum value is set based on the size of the block in memory (however, in our implementation, we set a fixed value of seven). In order to add a record in a block for storage, the record must first be hashed, as previously mentioned, and then undergo a verification process in which the hash and the digital signature of the record are verified for correctness across all nodes [3]. Our implementation does not utilize this verification process as we do not implement digital signatures, so we do not take this part of the process into account. Once verified, the record is added to the block, and a check is run to discover whether or not the block has reached the maximum number of transactions. If this is the case, a process known as consensus occurs, which produces a block hash that can be added to a chain [5]. This consensus process will be detailed further in Section 4.

Finally, chains are data structures that link the blocks themselves together. Each block, within its header, stores information about the hash of its parent block. Within a chain, this information is used to establish the order of the blocks— each block is linked to its predecessor based on the stored parent hashes [5]. By creating a chain of blocks in this manner, a secure and immutable data structure is developed in which one block cannot be altered or removed. Nodes run periodic checks on chains to verify that they have not been tampered with. If a block were to be manipulated, the hash of the block would change. This would cause the parent hash of the subsequent block to become incorrect, resulting in a break of the chain and making any attempt at altering data easily detectable. This tamper protection mechanism is an incredibly important part of blockchain and ensures that chains are secure and free of malicious attacks.

In our implementation, we establish a chain structure within a Merkle Tree, which will be detailed in section 3.

## 3&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Merkle Trees
Merkle Trees, also known as binary hash trees, are derivatives of binary tree data structures that store hashes within leaf nodes in order to maintain and verify the authenticity and integrity of data sets [7]. 

Like binary trees, Merkle Trees maintain their information within nodes. Each node in a Merkle Tree serves one of three purposes based on its placement. In figure 3.1, the nodes at the very bottom of the tree are known as leaves, and they have no children. Intermediate nodes, or nodes that have both a parent and children, are called branches, and they are represented in the second row of the tree figure. Finally, each Merkle Tree has a root node, which sits at the top of the tree and has no parent node associated with it [7]. Unlike binary trees, however, any node besides a leaf must have exactly two children. If a node has only a single child, it would be impossible to perform the hashing functionality, which will be detailed later in this section, and the tree would not be able to function.

Leaves store the hashes of records in the blockchain [7]. When a new record is received by the tree, the hashing function is run to generate the hash, which is then added to the tree. By storing hash values rather than the data of the transactions themselves, the information is unable to be retrieved, serving as a form of tamper protection. Leaves are stored in the Merkle Tree in the order in which they are received, and they are paired up as such under each branch node.

Branches are intermediary nodes, meaning, in the context of Merkle Trees, that they always have two child nodes. The branches directly above the leaf nodes in the Merkle Tree represent a concatenated hash of their two child nodes— that is, the hashes of the two child nodes are added together and then rehashed, which forms the hash that is stored within the branch. This methodology continues recursively up the tree up to the root of the tree. If a circumstance were to arise in which a branch node has only a single leaf, the leaf node will be duplicated up to the branch node and removed from the tree, effectively merging the two together and saving storage space in the process. [7]

Finally, the root of the tree, which has no parent node, is the first initial node added to the Merkle Tree. Once all of the branch and transaction hashing has occurred and been stored, the two children under the root are concatenated together and hashed. The result of this function is stored within the root hash. Since hashing works its way up from the bottom of the tree to the top, the root node is essentially a representation of the entire data set. This makes verification of transactions very simple. Rather than working down the tree to see if a transaction exists within the structure, the root hash can be checked for transaction consistency and transaction existence. In figure 3.1, querying for “h2” will return a hash that contains information about every other node in the tree, thereby proving that h2 must exist within the tree [7]. If for some reason the return value of the verification function is not consistent with this expected result, this would suggest that the Merkle Tree has been tampered with, and the chain would be considered invalid.

In our implementation, we maintain a Merkle Tree that stores hashes of records, and we also maintain a Merkle Tree within a chain that receives and stores hashes of blocks in leaves, which will be covered later in section 5. 

## 4&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Consensus
Consensus algorithms form the foundation of the process of adding blocks to chains. In blockchain consensus, all nodes within the network must agree when a block is received before it can be added to the blockchain. Blockchain typically utilizes one of two algorithms in order to do so: Proof of Work (PoW) and Proof of Stake (PoS). Broadly speaking, PoS uses less computational power and energy, but since it was released much more recently, it has not been as thoroughly tested and implemented within systems as PoW. PoW, first implemented in 2009 with the release of Bitcoin, has been thoroughly tested over the years. It is also much simpler to implement than PoS, though it does exert an enormous amount of computational power in order to run [8]. This section will go over the process and details of PoW, though both algorithms have valid uses and are widely used in cryptocurrency and other distributed systems.

Proof of Work, also known as mining in the context of blockchain, is a decentralized algorithmic mechanism that allows nodes in a blockchain network to verify blocks prior to adding them to a chain [8]. Nodes performing Proof of Work are known as miners. In PoW, miners compete to solve a computational puzzle in the form of a hash function in order to produce a valid block hash that meets a constraint known as a network target. In order to do so, miners utilize a value known as a nonce, or “number used once”, which is initialized to zero [8]. A miner generates a hash of the block using both the nonce and the block’s data. The nonce is incremented by one each time a hash is generated [8], meaning that the hash can be continuously regenerated without changing any information within the block, greatly reducing the chances of an attacker manipulating block data.

Each time a new hash is generated, it is compared against a value known as the network target. The target value is set based on the current computational power available in the blockchain network [8]. A miner must generate a hash which is less than or equal to this target value in order to solve the computational puzzle. If more computational resources are available on the network, the target value is increased [8]. Miners with more computational power can execute hashing functions more quickly, giving them an advantage over those with less power when it comes to mining and finding a valid hash. As a result, miners do not hold an equal stake in PoW, making it a competitive process. This competition, particularly in nodes with greater computational power, results in a much higher level of overall energy consumption compared to other consensus algorithms [8]. Once a hash has been found which meets the target requirements (greater than or equal to the target value), the block has been successfully validated and becomes available to add to a chain.

## 5&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Command Line Interface
In our implementation, our command-line interface provides users with two different options for interacting with the blockchain. When the program is initially run, a function is executed which initializes the connection of the nodes and instantiates a local chain for each connected node. In our implementation, we instantiate three chains, since we are simulating a network of three nodes. Once the nodes are connected, users are presented with an interface where they can enter either “1” or “2”.

Option 1 allows users to send a transaction to all nodes. Upon entering “1”, a user is presented with a prompt to enter the recipient’s information. In some blockchain implementations, such as cryptocurrency, this would be the ID of the wallet that a user is sending to. In our implementation, for simplicity, a user can enter any data since any transaction will be transmitted to every node via an RPC. After this, another prompt appears which allows the user to enter the data information that they want to send in the transaction. Any information can be entered here. Once this information has been entered, it is parsed by a reader generated from a function within the bufio Golang library and stored in two variables, “data” and “recipient”. These are then passed into a new transaction object, along with the address of the sending node and the current timestamp, generated by the using time.now() within the time library. The transaction is then distributed to all nodes using an RPC, and the CLI prints statements indicating whether or not the send and receipt was successful.

Option 2 allows users to view the hash of the current chain. For security reasons, users cannot directly access the chain itself, but they are able to view the hash of the root node in order to verify that their stored chain is consistent with the chains of other nodes connected to the network.

At the end of all of the prompts and information, the user is presented with another prompt indicating whether or not they want to continue. If they wish to continue and send or view information again, they would type “y” at this prompt. If they type “n”, the program will end and the node will exit the network.

## 6&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Implementation
The Blockchain is an intricate piece of code to write, split up into five smaller pieces to work with: Transaction, Block, Chain, PoW (Proof of Work), and Connection. 

Transactions populate Blocks in the form of a Transaction Merkle Tree called dataList. Transactions are added to the Merkle Tree, then stored within the block. Blocks, which are filled with transactions, are populated into a Block Merkle Tree, which is the blockchain. In order for the Blocks to be added into the blockchain, PoW is run to mine the block, which is added if a valid hash is found. Finally, Connection connects the peer nodes within the cluster to each other, allowing for chain consensus to occur as well as the mass sending of transactions and blocks to every node within the cluster. 

**Consensus**

The main purpose of blockchain is to create a decentralized, distributed ledger. Thus, consensus plays an important role in maintaining equivalent blockchains between all of the peers within the cluster. In order to do so, we use RPCs to send blocks to peers. The peers that receive the RPC will take the block data that was sent to them and construct a new block identical to the one that the node sent. An important note is that SendBlock() does not actually send an entire block data type; rather, it sends only the vital information to ensure that another node can reconstruct the block. This is because sending an entire block structure surpasses the byte limit that is allowed for an RPC. Sending a Transaction through SendTransaction() is similar to sending a Block: the node will receive only the vital information for it to accurately reconstruct the transaction and add it to its current block.
Transactions are sent through our CLI (command line interface), which was described in section 5. Whenever a user sends a transaction created through the CLI, AddTransaction()  is  executed.  This  deals with   all   the   logic   in   maintaining  chain  consensus with the other peers in the cluster. Say a transaction t is created by Node [A]. If the block b is not full, then [A] will add t to b. Note that b is a block that has not been added to the blockchain yet. Then, [A] will send t to all of its peer nodes. If b is 1 transaction away from max, then [A] will add t to b, making it full and executes AddBlock(), which mines the block and adds it to the blockchain. Then, [A] will SendBlock() b to all the peers it is connected to. Finally, [A] will initialize a new (empty) block that is not added to the chain, and send the empty block to its peers. 

Receiving RPCs, on the other hand, use functions that differ from those of the node that created the block or transaction, since the purpose of the peers is to verify that the block produces a valid hash before adding the block to its own blockchain, or to receive a transaction and add it to its own block without sending another RPC. For example, the receiver node does not use AddTransaction() because this function takes care of consensus and sends blocks and/or transactions to other nodes. If a peer who receives an RPC needs to add a transaction to its own block, it will use an internal function (which is also used in AddTransaction()) to add it. This is because we do not want nodes to send the same transaction back and forth. Adding a block is similar: instead of using AddBlock(), AddConsensusBlock() is used to verify that the data + nonce from block k, when recomputed in a different node, would result in the same hash as k.

## 7&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Limitations
As this is a simple implementation, there are limitations in terms of features. One limitation is that a node cannot join late because there is no system set in place to update the blockchain of latecomers. This could become a serious issue if the blockchain were to contain many blocks as the nodes would be out of sync in terms of stored information. To be clear, our implementation only works for nodes that are present at the time of the first transaction.

This drawback can be fixed with some RPCs to see which node has the longest blockchain. The node that joined late can send an RPC to all nodes that it is connected to within the network requesting the length of its own blockchain. Once all of the nodes respond, the node that joined late will contact the node with the longest blockchain and request an updated version. Another RPC can incrementally send all the blocks within the longest blockchain to the late node with SendBlock(), and the late node can build its own chain.

Another limitation within our implementation is the lack of digital signatures and security keys within transactions. In typical blockchain technology, digital signatures are an essential component of security, providing another layer of security as other nodes in the network are required to verify the signature of a transaction if a node initiates one. Digital signatures are created using a private key stored within a node, and then the transaction is verified by the other nodes using a public key accessible from the sending node. If for some reason the signatures do not match, the transaction is rejected. In our implementation, we did not implement this, though we do have other security mechanisms in place to ensure that transactions are valid, such as hashing.

## 8&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Conclusion
Despite the simplicity of our implementation, it serves as a foundation for understanding the underlying concepts of blockchain. However, it is important to note that this is a simplified version of blockchain technology and does not encompass all the features and security measures found in real-world implementations, such as digital signatures and the ability for late-joining nodes to synchronize their chains.

Overall, this project provides an introduction to blockchain and its core components, offering a basic understanding of how distributed ledger technology functions and its potential applications in various industries.

## References
[1] A Reuters Visual Guide Blockchain Explained
https://www.reuters.com/graphics/TECHNOLOGY-BLOCKCHAIN/010070P11GN/index.html   
[2] Building a Blockchain in Go PT:II - Proof of Work
https://dev.to/nheindev/building-a-blockchain-in-go-pt-ii-proof-of-work-eel  
[3] Di Pierro, Massimo. (2017). What Is the Blockchain?. Computing in Science & Engineering. 19. 92-95. 10.1109/MCSE.2017.3421554.    
[4] Introduction to Merkle Tree
https://www.geeksforgeeks.org/introduction-to-merkle-tree/       
[5] L. Ghiro et al., "A Blockchain Definition to Clarify its Role for the Internet of Things," 2021 19th Mediterranean Communication and Computer Networking Conference (MedComNet), Ibiza, Spain, 2021, pp. 1-8, doi: 10.1109/MedComNet52149.2021.9501280.       
[6] Merkle Tree
https://pkg.go.dev/github.com/cbergoon/merkletree         
[7] Merkle Tree in Blockchain: What it is and How it Works
https://www.investopedia.com/OHterms/m/merkle-tree.asp        
[8] What Is Proof of Work (PoW) in Blockchain? https://www.investopedia.com/terms/p/proof-work.aspk        
