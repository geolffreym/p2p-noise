package noise

import (
	"sync"
)

// [Socket] aliases for string.
type Socket string

// Bytes return a byte slice representation for socket.
func (s Socket) Bytes() []byte {
	return []byte(s)
}

// Bytes return a string representation for socket.
func (s Socket) String() string {
	return string(s)
}

// table assoc Socket with peer.
type table map[Socket]*peer

// Add new peer to table.
func (t table) Add(peer *peer) {
	t[peer.Socket()] = peer
}

// Get peer from table.
func (t table) Get(socket Socket) *peer {
	// exist socket related peer in table?
	if peer, ok := t[socket]; ok {
		return peer
	}

	return nil
}

// Remove peer from [Table].
func (t table) Remove(peer *peer) {
	delete(t, peer.Socket())
}

// router keep a hash table to associate [Socket] with peer.
// It is a hash table to associate Socket with Peers in a unstructured mesh topology.
type router struct {
	sync.RWMutex
	table table
}

func newRouter() *router {
	return &router{
		table: make(table),
	}
}

// Table return current routing Table.
func (r *router) Table() table { return r.table }

// Query return connection interface based on socket parameter.
func (r *router) Query(socket Socket) *peer {
	// Mutex for reading topics.
	// Do not write while topics are read.
	// Write Lock can’t be acquired until all Read Locks are released.
	// [RWMutex.Lock]: https://pkg.go.dev/sync#RWMutex.RLock
	r.RWMutex.RLock()
	defer r.RWMutex.RUnlock()
	// exist socket related peer?
	return r.table.Get(socket)
}

// Add create new Socket Peer association.
func (r *router) Add(peer *peer) {
	// Lock write/read table while add operation
	// A blocked Lock call excludes new readers from acquiring the lock.
	// [RWMutex.Lock]: https://pkg.go.dev/sync#RWMutex.Lock
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
func (r *router) Remove(peer *peer) {
	// Lock write/read table while add operation
	// A blocked Lock call excludes new readers from acquiring the lock.
	// [RWMutex.Lock]: https://pkg.go.dev/sync#RWMutex.Lock
	r.RWMutex.Lock()
	r.table.Remove(peer)
	r.RWMutex.Unlock()
}
