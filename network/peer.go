package network

import "net"

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

// PeerImp struct has a simplistic interface to describe a node in the network.
// Each peer has a socket address to identify itself and a connection interface to communicate with it.
// Each peer is routed and handled by route table.
type PeerImp struct {
	socket Socket   // IP and Port address for peer. https://en.wikipedia.org/wiki/Network_socket
	conn   net.Conn // Connection interface net.Conn to reach peer.
}

// Peer factory
func NewPeer(socket Socket, conn net.Conn) *PeerImp {
	return &PeerImp{
		conn:   conn,
		socket: socket,
	}
}

// Return peer connection interface
func (p *PeerImp) Connection() net.Conn { return p.conn }

// Return peer socket
func (p *PeerImp) Socket() Socket { return p.socket }

// Write buffered message over connection
func (p *PeerImp) Send(data []byte) (n int, err error) { return p.conn.Write(data) }

// Read buffered message from connection
func (p *PeerImp) Receive(buf []byte) (n int, err error) { return p.conn.Read(buf) }

// Close peer connection
func (p *PeerImp) Close() error { return p.conn.Close() }
