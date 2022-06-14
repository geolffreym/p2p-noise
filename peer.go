package noise

import (
	"bufio"
	"net"
)

// Peer struct has a simplistic interface to describe a peer in the network.
// Each Peer has a socket address to identify itself and a connection interface to communicate with it.
type Peer struct {
	socket   Socket            // IP and Port address for peer. https://en.wikipedia.org/wiki/Network_socket
	buffer   *bufio.ReadWriter // buffered IO. ref: https://golangdocs.com/bufio-package-golang
	net.Conn                   // embedded net.Conn to peer. ref: https://go.dev/doc/effective_go#embedding
}

// newBufferedIO creates a new buffered r/w IO
func newBufferedIO(conn net.Conn) *bufio.ReadWriter {
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
}

// Peer factory
func newPeer(socket Socket, conn net.Conn) *Peer {
	return &Peer{
		socket,
		newBufferedIO(conn),
		conn,
	}
}

// Return peer socket
func (p *Peer) Socket() Socket { return p.socket }

// Write buffered message over connection
func (p *Peer) Send(data []byte) (n int, err error) {
	return p.buffer.Write(data)
}

// Read buffered message from connection
func (p *Peer) Receive(buf []byte) (n int, err error) {
	return p.buffer.Read(buf)
}
