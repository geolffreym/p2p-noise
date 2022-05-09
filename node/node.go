package node

import (
	"github.com/geolffreym/p2p-noise/network"
	"github.com/geolffreym/p2p-noise/pubsub"
)

type Node struct {
	Done       chan bool
	Network    *network.Network
	subscriber *pubsub.Subscriber
}

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

func (n *Node) Listen(addr string) (*Node, error) {
	_, err := n.Network.Listen(addr)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (n *Node) Dial(addr string) (*Node, error) {
	_, err := n.Network.Dial(addr)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (n *Node) Broadcast(msg []byte) {
	for _, peer := range n.Network.Table() {
		go func(p *network.Peer) {
			p.Write(msg)
		}(peer)
	}
}

func (n *Node) Unicast(dest network.Socket, msg []byte) {
	if route, ok := n.Network.Table()[dest]; ok {
		route.Write(msg)
	}
}

func (n *Node) Observe(cb pubsub.Observer) {
	n.subscriber.Listen(cb)
}

// // Abstraction/alias for network event listener interface
// func (n *Node) AddListener(event network.Event) *Node {
// 	n.Network.AddEventListener(event, h)
// 	return n
// }

func (n *Node) Close() {
	n.Network.Close()
	n.Done <- true
}
