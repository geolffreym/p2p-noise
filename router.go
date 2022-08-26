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

// Table `keep` a socket:connection mapping.
type Table map[Socket]Peer

// Add new peer to table.
func (t Table) Add(peer Peer) {
	t[peer.Socket()] = peer
}

// Remove peer from table.
func (t Table) Remove(peer Peer) {
	delete(t, peer.Socket())
}

// router implements Router interface.
// It is a hash table to associate Socket with Peers in a unstructured mesh topology.
type router struct {
	sync.RWMutex
	table Table
}

func newRouter() *router {
	return &router{
		table: make(Table),
	}
}

// Table return current routing table.
func (r *router) Table() Table { return r.table }

// Query return connection interface based on socket parameter.
func (r *router) Query(socket Socket) Peer {
	// Mutex for reading topics.
	// Do not write while topics are read.
	// Write Lock can’t be acquired until all Read Locks are released.
	// ref: https://pkg.go.dev/sync#RWMutex.Lock
	r.RWMutex.RLock()
	defer r.RWMutex.RUnlock()

	// exist socket related peer?
	if peer, ok := r.table[socket]; ok {
		return peer
	}

	return nil
}

// Add create new socket connection association.
func (r *router) Add(peer Peer) {
	// Lock write/read table while add operation
	// A blocked Lock call excludes new readers from acquiring the lock.
	// ref: https://pkg.go.dev/sync#RWMutex.Lock
	r.RWMutex.Lock()
	r.table.Add(peer)
	r.RWMutex.Unlock()
}

// Len return the number of routed connections.
func (r *router) Len() uint8 {
	// 255 max peers len supported
	// uint8 is enough for routing peers len
	return uint8(len(r.table))
}

// Flush clean table and return total peers removed.
// This will be garbage collected eventually.
func (r *router) Flush() uint8 {
	size := r.Len()
	// nil its a valid type for mapping since its a reference type.
	// ref: https://github.com/go101/go101/wiki/About-the-terminology-%22reference-type%22-in-Go
	r.table = nil
	return size

}

// Remove removes a connection from router.
// It return recently removed peer.
func (r *router) Remove(peer Peer) {
	// Lock write/read table while add operation
	// A blocked Lock call excludes new readers from acquiring the lock.
	// ref: https://pkg.go.dev/sync#RWMutex.Lock
	r.RWMutex.Lock()
	r.table.Remove(peer)
	r.RWMutex.Unlock()
}
