package main

import (
	"crypto/rand"
	"encoding/hex"
)

// TxRef is the hex representation of the sha256 of the json payload of a transaction
type TxRef string

func NewTransaction(node *Node, tips []TxRef) Transaction {
	var bytes = make([]byte, 6)
	rand.Read(bytes)

	return Transaction{
		KeyChange: node.Key,
		nonce:     bytes,
		refs:      tips,
		Sig:       node.Key,
	}
}

type Transaction struct {
	// KeyChange signals a KeyChange for a node, aka a DID Doc update
	KeyChange string
	nonce     []byte
	// Origin == node.ID
	Origin    string
	refs      []TxRef
	Sig       string
}

func (tx Transaction) Reference() TxRef {
	return TxRef(hex.EncodeToString(tx.nonce))
}