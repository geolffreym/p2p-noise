package node

import (
	"github.com/geolffreym/p2p-noise/network"
)

// type Roles struct {
// 	requester net.Conn
// 	requested net.Conn
// }
type Node struct {
	Done    chan bool
	Network *network.Network
}

func New() *Node {
	return &Node{
		Network: network.New(),
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
	for _, route := range n.Network.Table() {
		route.Write(msg)
	}
}

func (n *Node) Unicast(dest network.Socket, msg []byte) {
	if route, ok := n.Network.Table()[dest]; ok {
		route.Write(msg)
	}
}

// Abstraction/alias for network event listener interface
func (n *Node) AddListener(event network.Event, h network.Handler) *Node {
	n.Network.AddEventListener(event, h)
	return n
}

func (n *Node) Close() {
	n.Network.Close()
	n.Done <- true
}
