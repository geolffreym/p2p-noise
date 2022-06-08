package network

import "sync"

// Aliases to handle idiomatic `Socket` type
type Socket string
type Table map[Socket]Peer

type Router interface {
	Len() int
	Table() Table
	Add(peer Peer) Peer
	Delete(peer Peer)
	Query(socket Socket) Peer
}

// Router hash table to associate Socket with Peers.
// eg. {127.0.0.1:4000: Peer}
type router struct {
	sync.RWMutex
	table Table
}

func NewRouter() Router {
	return &router{
		table: make(Table),
	}
}

func (r *router) Table() Table { return r.table }

// Add create new socket connection association
func (r *router) Add(peer Peer) Peer {
	// Lock write/read table while add operation
	// A blocked Lock call excludes new readers from acquiring the lock.
	r.RWMutex.Lock()
	defer r.RWMutex.Unlock()
	r.table[peer.Socket()] = peer
	return peer
}

// Return connection interface based on socket
func (r *router) Query(socket Socket) Peer {
	if peer, ok := r.table[socket]; ok {
		return peer
	}

	return nil
}

// Len return the number of connections
func (r *router) Len() int {
	return len(r.table)
}

// Delete removes a connection from router
func (r *router) Delete(peer Peer) {
	// Lock write/read table while delete operation
	// A blocked Lock call excludes new readers from acquiring the lock.
	r.RWMutex.Lock()
	defer r.RWMutex.Unlock()
	delete(r.table, peer.Socket())
}
