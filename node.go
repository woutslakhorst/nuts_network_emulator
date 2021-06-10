package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	rand2 "math/rand"
	"time"
)

type Node struct {
	// ID for display purposes
	ID          string
	// Key for signing TX. The Key is used for signature simulation => Key == Sig
	Key         string
	cfg         Config
	connections []*Connection
	network     *Network
	dag         *DAG
	vdr         VDR

	terminate   bool
	rounds		int
}

func newKey() string {
	bytes := make([]byte, 6)
	rand.Read(bytes)

	return hex.EncodeToString(bytes)
}

func NewNode(cfg Config, network *Network) *Node {
	return &Node{
		ID:      newKey(),
		network: network,
		cfg:     cfg,
		dag:     NewDag(),
		rounds:  cfg.Rounds,
		vdr:     NewVDR(),
	}
}

func (n *Node) connectAll() {
	for _, o := range n.network.nodes {
		if n.ID != o.ID {
			c := NewConnection(n, o)
			c.Start()
			n.connections = append(n.connections, c)
		}
	}
}

// Connect for receiving incoming connections
func (n *Node) Connect(connection *Connection) {
	go func() {
		for !n.terminate {
			select {
			case m := <-connection.Channel():
				n.process(m, connection.publisher)
			case <-time.After(10 * time.Millisecond):
			}
		}
	}()
}

func (n *Node) process(message Message, publisher *Node) {
	for _, t := range message.tips {
		n.processRef(t, publisher)
	}
}

func (n *Node) processRef(txRef TxRef, publisher *Node) {
	if !n.dag.Has(txRef) {
		//println(fmt.Sprintf("Node %s received new tx(%s) from %s", n.ID, txRef, publisher.ID))
		tx := publisher.GetTransaction(txRef)
		n.validate(tx)
		for _, p := range tx.refs {
			n.processRef(p, publisher)
		}
		n.dag.Offer(tx)
	}
}

func (n *Node) validate(tx *Transaction) {
	if !n.vdr.Matches(tx.Origin, tx.Sig) {
		println("ERROR incorrect SIG")
	}
	// todo: change timings/variance of published txs to simulate out of order processing
	n.vdr.Put(tx.Origin, tx.KeyChange)
}

func (n *Node) disconnectAll() {
	for _, c := range n.connections {
		c.Stop()
	}
}

func (n *Node) publish() {
	rate := 1.0
	if n.cfg.Rate != 0.0 {
		rate = n.cfg.Rate
	}
	variance := 0.1
	if n.cfg.Variance != 0.0 {
		variance = n.cfg.Variance
	}
	rate = ((rand2.Float64() * 2.0 * variance) - variance) * rate + rate

	go func() {
		for !n.terminate {
			time.Sleep(50 * time.Millisecond)
			n.sendTips()
		}
	}()

	go func() {
		for !n.terminate {
			time.Sleep(time.Duration(1000/rate) * time.Millisecond)
			n.NewTransaction()
			n.rounds--
			if n.rounds == 0 {
				time.Sleep(10 * time.Second)
				n.terminate = true
			}
		}
	}()

	go func() {
		for !n.terminate {
			time.Sleep(5 * time.Second)
			println(fmt.Sprintf("%s.tips (%d): %v", n.ID, n.cfg.Rounds - n.rounds, n.dag.Tips()))
		}
	}()
}

func (n *Node) sendTips() {
	for _, c := range n.connections {
		c.Send(Message{tips: n.dag.Tips()})
	}
}

func (n *Node) GetTransaction(txRef TxRef) *Transaction {
	return n.dag.GetTransaction(txRef)
}

func (n *Node) NewTransaction() {
	tips := n.dag.Tips()

	refs := make([]TxRef, len(tips))
	for i, e := range tips {
		refs[i] = e
	}
	if n.rounds % 100 == 0 {
		n.Key = newKey()
	}
	tx := NewTransaction(n, tips)

	n.dag.Offer(&tx)
}

func (n *Node) Start() error {
	// start connections
	n.connectAll()

	// start publisher
	n.publish()

	println(fmt.Sprintf("Started node %s", n.ID))
	return nil
}

func (n *Node) Stop() error {
	n.disconnectAll()
	n.terminate = true
	println(fmt.Sprintf("Stopped node %s", n.ID))
	return nil
}