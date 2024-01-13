package main

import (
	"fmt"
	"time"
)

// message in TokenRing
type Message struct {
	Data      string
	Recipient int
	TTL       int
}

// Nodering is node in ringtoken
type Nodering struct {
	id     int
	input  chan Message
	output chan Message
}

// TokenRing network
type networkTokenring struct {
	Nodes []*Nodering
}

// initializes n nodes in tokenring
func initialize(N int) *networkTokenring {
	network := &networkTokenring{}
	network.Nodes = make([]*Nodering, N)

	// initialize nodes, run goroutines
	for i := 0; i < N; i++ {
		network.Nodes[i] = &Nodering{
			id:     i,
			input:  make(chan Message),
			output: make(chan Message),
		}
		go network.runNode(network.Nodes[i])
	}
	for i := 0; i < N; i++ {
		network.Nodes[i].connect(network.Nodes[(i+1)%N], network.Nodes[(i-1+N)%N])
	}

	return network
}

// Connects nodes (input output)
func (n *Nodering) connect(next, prev *Nodering) {
	go func() {
		for message := range n.input {
			// send messsage to next node
			next.output <- message
		}
	}()

	go func() {
		for message := range prev.output {
			// receive message from previous node
			n.input <- message
		}
	}()
}

// node runs in the network token
func (network *networkTokenring) runNode(node *Nodering) {
	for message := range node.input {
		if message.Recipient == node.id {
			fmt.Printf("Message for node number %d : %s\n", node.id, message.Data)
		} else {
			// send message to next node
			node.output <- message
		}
	}
}

// sends a message to specific node in the network tokenring
func (network *networkTokenring) SendMessageToSpecificNode(data string, recipient int, ttl int) {
	message := Message{
		Data:      data,
		Recipient: recipient,
		TTL:       ttl,
	}
	// send message to first node
	network.Nodes[0].input <- message
}

func main() {

	var N, recipient, ttl int

	//input nodes n
	fmt.Print("Enter the number of nodes: ")
	fmt.Scan(&N)

	// initialize n nodes
	network := initialize(N)

	//input recipient
	fmt.Print("Enter recipient node (0 to N-1): ")
	fmt.Scan(&recipient)

	//input TTL
	fmt.Print("Enter TTL (time to live): ")
	fmt.Scan(&ttl)

	message := fmt.Sprintf("Hello,nodes %d!", recipient)

	// send message to specific node
	network.SendMessageToSpecificNode(message, recipient, ttl)

	// time delay message send
	time.Sleep(5 * time.Second)
}
