// Node describe a high level network interface.
// Each node is a Peer in the network and hold the needed methods to interact with other nodes.
package node

import (
	"context"
	"log"

	"github.com/geolffreym/p2p-noise/network"
	"github.com/geolffreym/p2p-noise/utils"
)

type NodeConnection interface {
	Listen(addr string) error
	Dial(addr string) error
	Close() error
}

type NodeSubscriber interface {
	Observe(cb network.Observer) context.CancelFunc
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
	messenger network.Messenger // Subscriber interface
	network   network.Network   // Network interface
}

// Build a ready to use subscriber interface and register default events for node
func NewNodeSubscriber(net network.Network) network.Messenger {
	subscriber := network.NewMessenger()
	net.Register(network.SELF_LISTENING, subscriber)
	net.Register(network.NEWPEER_DETECTED, subscriber)
	net.Register(network.MESSAGE_RECEIVED, subscriber)
	net.Register(network.CLOSED_CONNECTION, subscriber)
	net.Register(network.PEER_DISCONNECTED, subscriber)
	return subscriber
}

// Node factory
func NewNode() Node {

	// Register default events for node
	network := network.New()
	subscriber := NewNodeSubscriber(network)

	return &node{
		network:   network,
		messenger: subscriber,
	}
}

// Listen node in address and return error if it failed.
func (n *node) Listen(addr string) error {
	err := n.network.Listen(addr)
	if err != nil {
		return err
	}

	return nil
}

// Dial to a remote node and return error if it failed.
func (n *node) Dial(addr string) error {
	err := n.network.Dial(addr)
	if err != nil {
		return err
	}

	return nil
}

// Send messages to all the connected peers
func (n *node) Broadcast(msg []byte) {
	for _, peer := range n.network.Table() {
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
	if peer, ok := n.network.Table()[dest]; ok {
		_, err := peer.Send(msg)
		if err != nil {
			// TODO move message to errors
			log.Printf("error sending unicast message: %v", err)
		}
	}
}

// Use it to keep waiting for incoming notifications from the network.
func (n *node) Observe(cb network.Observer) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	n.messenger.Listen(ctx, cb)
	return cancel
}

// Close node connections and destroy node
func (n *node) Close() error {
	n.network.Close()         // Close network connection
	utils.Clear(&n.messenger) // Reset subscriber state
	return nil
}
