package noise

import (
	"sync"
)

type Socket string

// Table `keep` a socket:connection mapping
type Table map[Socket]*Peer

// TODO handle this from conf
// Max peers connected
const maxPeers = 255

// Add new peer to table
func (t Table) Add(peer *Peer) {
	t[peer.Socket()] = peer
}

// Remove peer from table
func (t Table) Remove(peer *Peer) {
	delete(t, peer.Socket())
}

// Router hash table to associate Socket with Peers.
// Unstructured mesh architecture
// eg. {127.0.0.1:4000: Peer}
type Router struct {
	sync.RWMutex
	table Table
}

func newRouter() *Router {
	return &Router{
		table: make(Table),
	}
}

// Table return current routing table
func (r *Router) Table() Table { return r.table }

// Return connection interface based on socket
func (r *Router) Query(socket Socket) *Peer {
	// Mutex for reading topics.
	// Do not write while topics are read.
	// Write Lock can’t be acquired until all Read Locks are released.
	// ref: https://pkg.go.dev/sync#RWMutex.Lock
	r.RWMutex.RLock()
	defer r.RWMutex.RUnlock()

	if peer, ok := r.table[socket]; ok {
		return peer
	}

	return nil
}

// Add create new socket connection association.
func (r *Router) Add(peer *Peer) {
	// Lock write/read table while add operation
	// A blocked Lock call excludes new readers from acquiring the lock.
	// ref: https://pkg.go.dev/sync#RWMutex.Lock
	if r.Len() > maxPeers {
		// TODO return error here
	}

	r.RWMutex.Lock()
	r.table.Add(peer)
	r.RWMutex.Unlock()
}

// Len return the number of connections
func (r *Router) Len() uint8 {
	// 255 max peers len supported
	// uint8 is enough for routing peers len
	return uint8(len(r.table))
}

// Remove removes a connection from router.
// It return recently removed peer.
func (r *Router) Remove(peer *Peer) {
	// Lock write/read table while add operation
	// A blocked Lock call excludes new readers from acquiring the lock.
	// ref: https://pkg.go.dev/sync#RWMutex.Lock
	if r.Len() == 0 {
		// TODO return error here
	}

	r.RWMutex.Lock()
	r.table.Remove(peer)
	r.RWMutex.Unlock()
}
