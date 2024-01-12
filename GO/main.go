package main

import (
	"fmt"
	"time"
)

// Token represents a message in the TokenRing
type Token struct {
	Data      string
	Recipient int
	TTL       int
}

// Node represents a node in the TokenRing
type Node struct {
	ID     int
	Input  chan Token
	Output chan Token
}

// TokenRing represents the TokenRing network
type TokenRing struct {
	Nodes []*Node
}

// Initializes a TokenRing with N nodes
func InitializeTokenRing(N int) *TokenRing {
	ring := &TokenRing{}
	ring.Nodes = make([]*Node, N)

	// Initialize nodes and launch goroutines
	for i := 0; i < N; i++ {
		ring.Nodes[i] = &Node{
			ID:     i,
			Input:  make(chan Token),
			Output: make(chan Token),
		}
		go ring.runNode(ring.Nodes[i])
	}

	// Connect nodes in a circular chain
	for i := 0; i < N; i++ {
		ring.Nodes[i].connect(ring.Nodes[(i+1)%N], ring.Nodes[(i-1+N)%N])
	}

	return ring
}

// Connects nodes with input and output channels
func (n *Node) connect(next, prev *Node) {
	go func() {
		for {
			select {
			case token := <-n.Input:
				// Forward token to the next node
				next.Output <- token
			}
		}
	}()

	go func() {
		for {
			select {
			case token := <-prev.Output:
				// Receive token from the previous node
				n.Input <- token
			}
		}
	}()
}

// Runs a node in the TokenRing
func (r *TokenRing) runNode(node *Node) {
	for {
		select {
		case token := <-node.Input:
			// Process the received token
			if token.Recipient == node.ID {
				fmt.Printf("Node %d received message: %s\n", node.ID, token.Data)
			} else {
				// Forward the token to the next node
				node.Output <- token
			}
		}
	}
}

// Sends a message to a specific node in the TokenRing
func (r *TokenRing) SendMessage(data string, recipient int, ttl int) {
	token := Token{
		Data:      data,
		Recipient: recipient,
		TTL:       ttl,
	}

	// Send the token to the first node in the TokenRing
	r.Nodes[0].Input <- token
}

func main() {
	N := 5 // Number of nodes in the TokenRing

	// Initialize the TokenRing with N nodes
	ring := InitializeTokenRing(N)

	// Send a message from the main thread to a specific node
	ring.SendMessage("Hello, Node 3!", 4, 2)

	// Allow some time for the message to propagate through the TokenRing
	time.Sleep(5 * time.Second)
}
