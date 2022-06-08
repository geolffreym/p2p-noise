package network

import (
	"net"
)

type PeerConnection interface {
	// Return peer connection interface
	Connection() net.Conn
	// Return peer socket
	Socket() Socket
	// Close peer connection
	Close() error
}

type PeerStreamer interface {
	// Write buffered message over connection
	Send(data []byte) (n int, err error)
	// Read buffered message from connection
	Receive(buf []byte) (n int, err error)
}

type Peer interface {
	PeerConnection
	PeerStreamer
}

// peer struct has a simplistic interface to describe a node in the network.
// Each peer has a socket address to identify itself and a connection interface to communicate with it.
// Each peer is routed and handled by route table.
type peer struct {
	socket Socket   // IP and Port address for peer. https://en.wikipedia.org/wiki/Network_socket
	conn   net.Conn // Connection interface net.Conn to reach peer.
}

// Peer factory
func NewPeer(socket Socket, conn net.Conn) Peer {
	return &peer{
		conn:   conn,
		socket: socket,
	}
}

// Return peer connection interface
func (p *peer) Connection() net.Conn { return p.conn }

// Return peer socket
func (p *peer) Socket() Socket { return p.socket }

// Write buffered message over connection
func (p *peer) Send(data []byte) (n int, err error) { return p.conn.Write(data) }

// Read buffered message from connection
func (p *peer) Receive(buf []byte) (n int, err error) {
	// p.conn.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
	return p.conn.Read(buf)
}

// Close peer connection
func (p *peer) Close() error { return p.conn.Close() }
