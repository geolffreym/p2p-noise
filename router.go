package noise

import (
	"sync"
	"sync/atomic"
)

// router keeps a hash table to associate IDs with peers.
// It implements an unstructured mesh topology.
// Unstructured P2P topologies do not attempt to organize all peers into a single, structured topology.
// Rather, each peer attempts to keep a "sensible" set of other peers in its routing table.
type router struct {
	sync.Map // embed map
	counter  uint32
}

func newRouter() *router {
	return &router{counter: 0}
}

// Table return fan out channel with routed peers.
func (r *router) Table() <-chan *peer {
	// buffered channel
	ch := make(chan *peer, r.Len())
	// ref: https://pkg.go.dev/sync#Map.Range
	// generate valid peers from table
	r.Range(func(_, value any) bool {
		if p, ok := value.(*peer); ok {
			ch <- p
		}
		// keep running until finish sequence
		// If f returns false, range stops the iteration.
		return true
	})

	close(ch)
	return ch
}

// Query return connection interface based on socket parameter.
// In-band error returned. This return value may be an error, or a boolean when no explanation is needed.
// refer:  https://go.dev/wiki/CodeReviewComments
func (r *router) Query(id ID) (*peer, bool) {

	p, exists := r.Load(id)
	peer, ok := p.(*peer)

	if !exists || !ok {
		return nil, false
	}

	return peer, true
}

// Add forward method to internal sync.Map store for peer.
func (r *router) Add(peer *peer) {
	atomic.AddUint32(&r.counter, 1)
	r.Store(peer.ID(), peer)
}

// Len return the number of routed connections.
func (r *router) Len() uint8 {
	return uint8(atomic.LoadUint32(&r.counter))
}

// Remove forward method to internal sync.Map to delete a connection from router.
// It return recently removed peer.
func (r *router) Remove(peer *peer) {
	// ref: https://github.com/golang/go/blob/509ee7064207cc9c8ac81bc76f182a5fbb877e9b/src/sync/atomic/doc.go#L96
	atomic.AddUint32(&r.counter, ^uint32(0))
	r.Delete(peer.ID())
}
