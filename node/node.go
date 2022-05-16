// Node describe a high level network interface.
// Each node is a Peer in the network and hold the needed methods to interact with other nodes.
package node

import (
	"log"

	"github.com/geolffreym/p2p-noise/network"
	"github.com/geolffreym/p2p-noise/utils"
)

type Node interface {
	Observe(cb network.Observer)
	Listen(addr string) (*Node, error)
	Dial(addr string) (*Node, error)
	Close()
}

// Node implementation.
type NodeImp struct {
	Sentinel   chan bool           // Hangs while waiting for channel closed and stop node process.
	Network    *network.Network    // Network interface
	subscriber *network.Subscriber // Subscriber interface
}

// Build a ready to use subscriber interface and register default events for node
func NodeSubscriber(n *network.Network) *network.Subscriber {
	subscriber := network.NewSubscriber()
	n.Events.Register(network.SELF_LISTENING, subscriber)
	n.Events.Register(network.NEWPEER_DETECTED, subscriber)
	n.Events.Register(network.MESSAGE_RECEIVED, subscriber)
	n.Events.Register(network.CLOSED_CONNECTION, subscriber)
	return subscriber
}

// Node factory
func NewNode() *NodeImp {

	// Register default events for node
	network := network.New()
	subscriber := NodeSubscriber(network)

	return &NodeImp{
		Network:    network,
		Sentinel:   make(chan bool),
		subscriber: subscriber,
	}
}

// Listen node in address and return error if it failed.
func (n *NodeImp) Listen(addr string) (*NodeImp, error) {
	_, err := n.Network.Listen(addr)
	if err != nil {
		return nil, err
	}

	return n, nil
}

// Dial to a remote node and return error if it failed.
func (n *NodeImp) Dial(addr string) (*NodeImp, error) {
	_, err := n.Network.Dial(addr)
	if err != nil {
		return nil, err
	}

	return n, nil
}

// Send messages to all the connected peers
func (n *NodeImp) Broadcast(msg []byte) {
	for _, peer := range n.Network.Table() {
		go func(n *NodeImp, p network.Peer) {
			_, err := p.Send(msg)
			if err != nil {
				// TODO move message to errors
				log.Printf("error sending broadcast message: %v", err)
			}

		}(n, peer)
	}
}

// Send a message to a specific peer
func (n *NodeImp) Unicast(dest network.Socket, msg []byte) {
	if peer, ok := n.Network.Table()[dest]; ok {
		_, err := peer.Send(msg)
		if err != nil {
			// TODO move message to errors
			log.Printf("error sending unicast message: %v", err)
		}
	}
}

// Use it to keep waiting for incoming notifications from the network.
func (n *NodeImp) Observe(cb network.Observer) {
	n.subscriber.Listen(cb)
}

// Close node connections and destroy node
func (n *NodeImp) Close() {
	n.Network.Close()          // Close network connection
	utils.Clear(&n.subscriber) // Reset subscriber state
	close(n.Sentinel)          // Stop node
}
