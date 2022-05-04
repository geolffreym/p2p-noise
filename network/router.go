package network

import (
	"net"
)

// Alias for string
type Socket string

type Peer struct {
	socket Socket
	conn   net.Conn
}

func (r *Peer) Connection() net.Conn                 { return r.conn }
func (r *Peer) Socket() Socket                       { return r.socket }
func (r *Peer) Write(data []byte) (n int, err error) { return r.conn.Write(data) }
func (r *Peer) Read(buf []byte) (n int, err error)   { return r.conn.Read(buf) }
func (r *Peer) Close() error                         { return r.conn.Close() }

type Router map[Socket]*Peer

func (r Router) Add(socket Socket, route *Peer) {
	r[socket] = route
}

func (r Router) Len() int {
	return len(r)
}

func (r Router) Delete(socket Socket) {
	delete(r, socket)
}
