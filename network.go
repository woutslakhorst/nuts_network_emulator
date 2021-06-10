package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Config struct {
	// Nodes is the number of nodes in the network
	Nodes int
	// Rounds is the number of Txs each node will publish, -1 for indefinite
	Rounds int
	// Rate is the number of Txs per second, 1.0 if not defined
	Rate float64
	// Variance is how much the rate may differ. 0.1 == 10%. 0.1 if not defined
	Variance float64
}

type Network struct {
	nodes []*Node
}

func NewNetwork(cfg Config) *Network {
	n := Network{}
	n.nodes = make([]*Node, cfg.Nodes)
	for i := 0; i < cfg.Nodes; i++ {
		n.nodes[i] = NewNode(cfg, &n)
	}

	return &n
}

func (n *Network) Start() error {
	// listen for signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// start nodes
	for _, n := range n.nodes {
		if err := n.Start(); err != nil {
			return err
		}
	}

	// blocking loop
	done := false
	for !done {
		select {
			case _ = <- sigs:
				done = true
			case _ = <- time.After(5 * time.Second):

		}
	}

	// stop all nodes
	for _, n := range n.nodes {
		if err := n.Stop(); err != nil {
			return err
		}
	}

	return nil
}
