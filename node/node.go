// Node describe a high level network interface.
// Each node is a Peer in the network and hold the needed methods to interact with other nodes.
package node

import (
	"github.com/geolffreym/p2p-noise/network"
	"github.com/geolffreym/p2p-noise/pubsub"
)

type Node interface {
	Observe(cb pubsub.Observer)
	Listen(addr string) (*Node, error)
	Dial(addr string) (*Node, error)
	Close()
}

// Node implementation.
type NodeImp struct {
	Sentinel   chan bool          // Hangs while waiting for channel closed and stop node process.
	Network    *network.Network   // Network interface
	subscriber *pubsub.Subscriber // Subscriber interface
}

// Node factory
func NewNode() *NodeImp {

	// Register default events for node
	network := network.New()
	subscriber := pubsub.NewSubscriber()
	network.Events.Register(pubsub.NEWPEER_DETECTED, subscriber)
	network.Events.Register(pubsub.SELF_LISTENING, subscriber)
	network.Events.Register(pubsub.MESSAGE_RECEIVED, subscriber)

	return &NodeImp{
		Network:    network,
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
		go func(p *network.Peer) {
			p.Write(msg)
		}(peer)
	}
}

// Send a message to a specific peer
func (n *NodeImp) Unicast(dest network.Socket, msg []byte) {
	if route, ok := n.Network.Table()[dest]; ok {
		route.Write(msg)
	}
}

// Use it to keep waiting for incoming notifications from the network.
func (n *NodeImp) Observe(cb pubsub.Observer) {
	n.subscriber.Listen(cb)
}

// Close node connections
func (n *NodeImp) Close() {
	n.Network.Close()
	close(n.Sentinel)
}
