package network

import "net"

// Node peer definition
type Peer struct {
	socket Socket   // IP and Port address for peer. https://en.wikipedia.org/wiki/Network_socket
	conn   net.Conn // Connection interface net.Conn to reach peer.
}

// Return peer connection interface
func (r *Peer) Connection() net.Conn { return r.conn }

// Return peer socket
func (r *Peer) Socket() Socket { return r.socket }

// Write buffered message over connection
func (r *Peer) Write(data []byte) (n int, err error) { return r.conn.Write(data) }

// Read buffered message from connection
func (r *Peer) Read(buf []byte) (n int, err error) { return r.conn.Read(buf) }

// Close peer connection
func (r *Peer) Close() error { return r.conn.Close() }
