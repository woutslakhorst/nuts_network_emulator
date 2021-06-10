package main

import (
	"fmt"
	"time"
)

type Connection struct {
	publisher *Node
	consumer  *Node
	connected bool
	terminate bool
	inChannel chan Message
	outChannel chan Message
}

func NewConnection(publisher *Node, consumer *Node) *Connection {
	return &Connection {
		publisher: publisher,
		consumer: consumer,
		inChannel: make(chan Message),
		outChannel: make(chan Message),
	}
}

func (c *Connection) Start() {

	go func() {
		for !c.terminate {
			if !c.connected {
				c.consumer.Connect(c)
				c.connected = true
				println(fmt.Sprintf("Node %s connected to node %s", c.publisher.ID, c.consumer.ID))
			} else {
				select {
				case m := <-c.inChannel:
					c.outChannel <- m
				case <-time.After(10 * time.Millisecond):
				}
			}
		}
		println(fmt.Sprintf("Node %s terminated from node %s", c.publisher.ID, c.consumer.ID))
	}()
}

func (c *Connection) Send(message Message) {
	c.inChannel <- message
}

func (c *Connection) Channel() chan Message {
	return c.outChannel
}

func (c *Connection) Stop() {
	c.terminate = true
}