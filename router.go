package noise

import (
	"sync"
)

// Handle socket string logic
// ip:port eg. 127.0.0.1:8000
type Socket string

func (s Socket) Bytes() []byte {
	return []byte(s)
}

func (s Socket) String() string {
	return string(s)
}

// table `keep` a socket:connection mapping.
type table map[Socket]*peer

// Add new peer to table.
func (t table) Add(peer *peer) {
	t[peer.Socket()] = peer
}

// Remove peer from table.
func (t table) Remove(peer *peer) {
	delete(t, peer.Socket())
}

// router hash table to associate Socket with Peers.
// Unstructured mesh architecture.
// eg. {127.0.0.1:4000: Peer}
type router struct {
	sync.RWMutex
	table table
}

func newRouter() *router {
	return &router{
		table: make(table),
	}
}

// Table return current routing table.
func (r *router) Table() table { return r.table }

// Return connection interface based on socket.
func (r *router) Query(socket Socket) *peer {
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
func (r *router) Add(peer *peer) {
	// Lock write/read table while add operation
	// A blocked Lock call excludes new readers from acquiring the lock.
	// ref: https://pkg.go.dev/sync#RWMutex.Lock
	r.RWMutex.Lock()
	r.table.Add(peer)
	r.RWMutex.Unlock()
}

// Len return the number of connections.
func (r *router) Len() uint8 {
	// 255 max peers len supported
	// uint8 is enough for routing peers len
	return uint8(len(r.table))
}

// Remove removes a connection from router.
// It return recently removed peer.
func (r *router) Remove(peer *peer) {
	// Lock write/read table while add operation
	// A blocked Lock call excludes new readers from acquiring the lock.
	// ref: https://pkg.go.dev/sync#RWMutex.Lock
	r.RWMutex.Lock()
	r.table.Remove(peer)
	r.RWMutex.Unlock()
}
