package noise

import (
	"bufio"
	"net"
)

// Peer struct has a simplistic interface to describe a peer in the network.
// Each Peer has a socket address to identify itself and a connection interface to communicate with it.
type Peer struct {
	net.Conn                   // embedded net.Conn to peer. ref: https://go.dev/doc/effective_go#embedding
	socket   Socket            // IP and Port address for peer. https://en.wikipedia.org/wiki/Network_socket
	buffer   *bufio.ReadWriter // buffered IO. why? ref: https://golangdocs.com/bufio-package-golang
}

// newBufferedIO creates a new buffered r/w IO.
func newBufferedIO(conn net.Conn) *bufio.ReadWriter {
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
}

func newPeer(socket Socket, conn net.Conn) *Peer {
	// Go does not provide the typical, type-driven notion of subclassing,
	// but it does have the ability to “borrow” pieces of an implementation by embedding types within a struct or interface.
	return &Peer{
		conn,
		socket,
		newBufferedIO(conn),
	}
}

// Return peer socket.
func (p *Peer) Socket() Socket { return p.socket }

// Write buffered message over connection.
func (p *Peer) Write(data []byte) (n int, err error) {
	// This forwarding method is needed to handle ambiguous method names
	return p.buffer.Write(data)
}

// Read buffered message from connection.
func (p *Peer) Read(buf []byte) (n int, err error) {
	// This forwarding method is needed to handle ambiguous method names
	return p.buffer.Read(buf)
}
