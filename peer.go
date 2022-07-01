package noise

import (
	"net"
)

// Peer struct has a simplistic interface to describe a peer in the network.
// Each Peer has a socket address to identify itself and a connection interface to communicate with it.
type Peer struct {
	net.Conn        // embedded net.Conn to peer. ref: https://go.dev/doc/effective_go#embedding
	socket   Socket // IP and Port address for peer. https://en.wikipedia.org/wiki/Network_socket
}

func newPeer(socket Socket, conn net.Conn) *Peer {
	// Go does not provide the typical, type-driven notion of sub-classing,
	// but it does have the ability to “borrow” pieces of an implementation by embedding types within a struct or interface.
	return &Peer{
		conn,
		socket,
	}
}

// Return peer socket.
func (p *Peer) Socket() Socket { return p.socket }
