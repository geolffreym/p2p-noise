package noise

import (
	"sync"
	"sync/atomic"
)

// router keep a hash table to associate ID with peer.
// It implements a unstructured mesh topology.
type router struct {
	sync.Map
	counter uint32
}

func newRouter() *router {
	return &router{counter: 0}
}

// Table return fan out channel with routed peers.
func (r *router) Table() <-chan *peer {
	ch := make(chan *peer, r.Len())
	// ref: https://pkg.go.dev/sync#Map.Range
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
func (r *router) Query(id ID) *peer {
	// exist socket related peer?
	p, exists := r.Load(id)
	peer, ok := p.(*peer)

	if !exists || !ok {
		return nil
	}

	return peer
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
