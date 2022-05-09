package network

import "net"

type Peer struct {
	socket Socket
	conn   net.Conn
}

func (r *Peer) Connection() net.Conn                 { return r.conn }
func (r *Peer) Socket() Socket                       { return r.socket }
func (r *Peer) Write(data []byte) (n int, err error) { return r.conn.Write(data) }
func (r *Peer) Read(buf []byte) (n int, err error)   { return r.conn.Read(buf) }
func (r *Peer) Close() error                         { return r.conn.Close() }
