package node

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/multiformats/go-multiaddr"
)

type Node struct {
	Host host.Host
	Room protocol.ID
	ctx  context.Context
}

func New(room string, port int) *Node {
	// 0.0.0.0 will listen on any interface device.
	localhostMAddress := fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port)
	multiAddress, _ := multiaddr.NewMultiaddr(localhostMAddress)

	// Create a libp2p host/peer
	host, _ := libp2p.New(libp2p.ListenAddrs(multiAddress))

	return &Node{
		Host: host,
		Room: protocol.ID(room),
		ctx:  context.Background(),
	}

}

func (p *Node) ID() peer.ID { return p.Host.ID() }

// Add stream handler to handle in/out messages
func (p *Node) AddStreamHandler(handler network.StreamHandler) {
	p.Host.SetStreamHandler(p.Room, handler)
}

// Create stream to specific peer
func (p *Node) StreamToPeer(peer peer.ID) (network.Stream, error) {
	return p.Host.NewStream(p.ctx, peer, p.Room)
}

// Bootstrapping DHT node
func (p *Node) Bootstrap() {

}
