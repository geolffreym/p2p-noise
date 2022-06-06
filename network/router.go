package network

import "sync"

// Aliases to handle idiomatic `Socket` type
type Socket string

// Router hash table to associate Socket with Peers.
// eg. {127.0.0.1:4000: Peer}

type Table map[Socket]Peer

type router struct {
	sync.RWMutex
	table Table
}

type Router interface {
	Table() Table
	Add(socket Socket, peer Peer)
	Delete(socket Socket)
	Len() int
}

func NewRouter() Router {
	return &router{
		table: make(Table),
	}
}

func (r *router) Table() Table {
	return r.table
}

// Add create new socket connection association
func (r *router) Add(socket Socket, peer Peer) {
	// Lock write/read table while routing operation
	// A blocked Lock call excludes new readers from acquiring the lock.
	r.RWMutex.Lock()
	defer r.RWMutex.Unlock()
	r.table[socket] = peer
}

// Len return the number of connections
func (r *router) Len() int {
	return len(r.table)
}

// Delete removes a connection from router
func (r *router) Delete(socket Socket) {
	// Lock write/read table while removing route
	// A blocked Lock call excludes new readers from acquiring the lock.
	r.RWMutex.Lock()
	defer r.RWMutex.Unlock()
	delete(r.table, socket)
}
