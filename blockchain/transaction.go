package blockchain

import (
	"bytes"
	"crypto/sha256"
	"reflect"

	"github.com/cbergoon/merkletree"
)

// sender and recipient in cryptocurrency are referring to the keys of the wallets
// in our implementation, we will use the ip address of the sender and recipient. This is only used for metadata for hashing purposes
// not actually sending anything to the recipient (no actual cryptocurrency)
// in real-life cryptocurrency, recipient does not have to be online to receive
// sender		ip address of the sender
// recipient	ip address of the recipient
// timestamp	time when transaction is created
// data			message that sender wants to send to recipient. In real-life cryptocurrency,
//
//	this is the amount of cryptocurrency that the sender wants to send to the recipient
type Transaction struct {
	Sender    []byte
	Recipient []byte
	Timestamp int64
	Data      []byte
}

// Turns everything in the Transaction struct into a byte array
func (transaction *Transaction) TransactionDataToBytes() []byte {
	data := bytes.Join(
		[][]byte{
			transaction.Sender,
			transaction.Recipient,
			ToHex(int64(transaction.Timestamp)),
			transaction.Data,
		},
		[]byte{},
	)

	return data
}

// Calculates the hash of the transaction
// Implements the merkletree.Content interface
func (transaction Transaction) CalculateHash() ([]byte, error) {
	hash := sha256.New()
	data := transaction.TransactionDataToBytes()

	if _, err := hash.Write(data); err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}

// Implements the merkletree.Content interface
func (transaction Transaction) Equals(other merkletree.Content) (bool, error) {
	ifEquals := false

	if reflect.DeepEqual(transaction, other.(Transaction)) {
		ifEquals = true
	}

	return ifEquals, nil
}
