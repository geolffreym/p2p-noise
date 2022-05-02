package network

import (
	"bufio"
	"net"
)

// Alias for string
type socket string
type router map[socket]*Route
type handler func(*Route)

func (r router) Add(socket socket, route *Route) {
	r[socket] = route
}

func (r router) Connected() int {
	return len(r)
}

func (r router) Delete(socket socket) {
	delete(r, socket)
}

type Route struct {
	socket socket
	conn   net.Conn
	stream *bufio.ReadWriter
}

func (r *Route) Stream() *bufio.ReadWriter { return r.stream }
func (r *Route) Connection() net.Conn      { return r.conn }
func (r *Route) Socket() socket            { return r.socket }
func (r *Route) Close()                    { r.conn.Close() }
