package noise

import (
	"sync"
)

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

func (r *Router) Table() Table { return r.table }

// Return connection interface based on socket
func (r *Router) Query(socket Socket) *Peer {
	// Mutex for reading topics.
	// Do not write while topics are read.
	// Write Lock can’t be acquired until all Read Locks are released.
	r.RWMutex.RLock()
	defer r.RWMutex.RUnlock()

	if peer, ok := r.table[socket]; ok {
		return peer
	}

	return nil
}

// Add create new socket connection association
func (r *Router) Add(peer *Peer) {
	// Lock write/read table while add operation
	// A blocked Lock call excludes new readers from acquiring the lock.
	r.RWMutex.Lock()
	defer r.RWMutex.Unlock()
	r.table[peer.Socket()] = peer
}

// Len return the number of connections
func (r *Router) Len() int {
	return len(r.table)
}

// Delete removes a connection from router
func (r *Router) Delete(peer *Peer) {
	// Lock write/read table while delete operation
	// A blocked Lock call excludes new readers from acquiring the lock.
	r.RWMutex.Lock()
	defer r.RWMutex.Unlock()
	delete(r.table, peer.Socket())
}
