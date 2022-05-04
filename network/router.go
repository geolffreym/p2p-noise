package network

import (
	"net"
)

// Alias for string
type Socket string

type Route struct {
	socket Socket
	conn   net.Conn
}

func (r *Route) Connection() net.Conn                 { return r.conn }
func (r *Route) Socket() Socket                       { return r.socket }
func (r *Route) Write(data []byte) (n int, err error) { return r.conn.Write(data) }
func (r *Route) Read(buf []byte) (n int, err error)   { return r.conn.Read(buf) }
func (r *Route) Close() error                         { return r.conn.Close() }

type Router map[Socket]*Route

func (r Router) Add(socket Socket, route *Route) {
	r[socket] = route
}

func (r Router) Len() int {
	return len(r)
}

func (r Router) Delete(socket Socket) {
	delete(r, socket)
}
