package node

import (
	"github.com/geolffreym/p2p-noise/network"
	"github.com/geolffreym/p2p-noise/pubsub"
)

// Node describe a high level network interface.
// Each node is a Peer in the network and hold the needed methods to interact with other nodes.
type Node struct {
	Done       chan bool          // Done hangs while waiting for be closed
	Network    *network.Network   // Network interface
	subscriber *pubsub.Subscriber // Subscriber interface
}

// Node factory
func New() *Node {

	// Register default events for node
	network := network.New()
	subscriber := pubsub.NewSubscriber()
	network.Events.Register(pubsub.NEWPEER_DETECTED, subscriber)
	network.Events.Register(pubsub.SELF_LISTENING, subscriber)
	network.Events.Register(pubsub.MESSAGE_RECEIVED, subscriber)

	return &Node{
		Network:    network,
		subscriber: subscriber,
	}
}

// Listen node in address
func (n *Node) Listen(addr string) (*Node, error) {
	_, err := n.Network.Listen(addr)
	if err != nil {
		return nil, err
	}

	return n, nil
}

// Dial to a remote node
func (n *Node) Dial(addr string) (*Node, error) {
	_, err := n.Network.Dial(addr)
	if err != nil {
		return nil, err
	}

	return n, nil
}

// Send messages to all the connected peers
func (n *Node) Broadcast(msg []byte) {
	for _, peer := range n.Network.Table() {
		go func(p *network.Peer) {
			p.Write(msg)
		}(peer)
	}
}

// Send a message to a specific peer
func (n *Node) Unicast(dest network.Socket, msg []byte) {
	if route, ok := n.Network.Table()[dest]; ok {
		route.Write(msg)
	}
}

// Use it to keep waiting for incoming notifications from the network.
func (n *Node) Observe(cb pubsub.Observer) {
	n.subscriber.Listen(cb)
}

// Close node connections
func (n *Node) Close() {
	n.Network.Close()
	close(n.Done)
}
