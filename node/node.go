package node

import (
	"github.com/geolffreym/p2p-network/network"
)

// type Roles struct {
// 	requester net.Conn
// 	requested net.Conn
// }
type Node struct {
	Addr    string
	Done    chan bool
	Network *network.Network
}

func New(addr string) *Node {
	return &Node{
		Addr:    addr,
		Network: network.New(addr),
	}
}

func (n *Node) Listen() (*Node, error) {
	conn, err := n.Network.Listen()
	if err != nil {
		return nil, err
	}

	conn.Bind() // waiting for peers
	return n, nil
}

func (n *Node) Dial(addr string) (*Node, error) {
	_, err := n.Network.Dial(addr)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (n *Node) Close() {
	n.Network.Close()
	n.Done <- true
}
