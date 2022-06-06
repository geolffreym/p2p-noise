// Node describe a high level network interface.
// Each node is a Peer in the network and hold the needed methods to interact with other nodes.
package node

import (
	"log"

	"github.com/geolffreym/p2p-noise/network"
	"github.com/geolffreym/p2p-noise/utils"
)

type NodeConnection interface {
	Listen(addr string) (Node, error)
	Dial(addr string) (Node, error)
	Close()
}

type NodeSubscriber interface {
	Observe(cb network.Observer)
	Wait()
}

type NodeEmitter interface {
	Broadcast(msg []byte)
	Unicast(dest network.Socket, msg []byte)
}

type Node interface {
	NodeConnection
	NodeSubscriber
	NodeEmitter
}

// Node implementation.
type node struct {
	sentinel   chan bool          // Hangs while waiting for channel closed and stop node process.
	subscriber network.Subscriber // Subscriber interface
	Network    network.Network    // Network interface
}

// Build a ready to use subscriber interface and register default events for node
func NewNodeSubscriber(net network.Network) network.Subscriber {
	subscriber := network.NewSubscriber()
	net.Register(network.SELF_LISTENING, subscriber)
	net.Register(network.NEWPEER_DETECTED, subscriber)
	net.Register(network.MESSAGE_RECEIVED, subscriber)
	net.Register(network.CLOSED_CONNECTION, subscriber)
	return subscriber
}

// Node factory
func NewNode() Node {

	// Register default events for node
	network := network.New()
	subscriber := NewNodeSubscriber(network)

	return &node{
		Network:    network,
		sentinel:   make(chan bool),
		subscriber: subscriber,
	}
}

// Listen node in address and return error if it failed.
func (n *node) Listen(addr string) (Node, error) {
	_, err := n.Network.Listen(addr)
	if err != nil {
		return nil, err
	}

	return n, nil
}

// Dial to a remote node and return error if it failed.
func (n *node) Dial(addr string) (Node, error) {
	_, err := n.Network.Dial(addr)
	if err != nil {
		return nil, err
	}

	return n, nil
}

// Send messages to all the connected peers
func (n *node) Broadcast(msg []byte) {
	for _, peer := range n.Network.Router().Table() {
		go func(n *node, p network.Peer) {
			_, err := p.Send(msg)
			if err != nil {
				// TODO move message to errors
				log.Printf("error sending broadcast message: %v", err)
			}

		}(n, peer)
	}
}

// Send a message to a specific peer
func (n *node) Unicast(dest network.Socket, msg []byte) {
	if peer, ok := n.Network.Router().Table()[dest]; ok {
		_, err := peer.Send(msg)
		if err != nil {
			// TODO move message to errors
			log.Printf("error sending unicast message: %v", err)
		}
	}
}

// Use it to keep waiting for incoming notifications from the network.
func (n *node) Observe(cb network.Observer) {
	n.subscriber.Listen(cb)
}

// Locked channel until the network get closed
func (n *node) Wait() {
	<-n.sentinel
}

// Close node connections and destroy node
func (n *node) Close() {
	n.Network.Close()          // Close network connection
	utils.Clear(&n.subscriber) // Reset subscriber state
	close(n.sentinel)          // Stop node
}
