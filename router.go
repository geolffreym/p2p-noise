package noise

import "sync"

// router keep a hash table to associate ID with peer.
// It implements a unstructured mesh topology.
type router struct {
	table   sync.Map
	counter uint8
}

func newRouter() *router {
	return &router{}
}

// Table return current routing table.
func (r *router) Table() <-chan *peer {
	ch := make(chan *peer, r.counter)
	r.table.Range(func(key, value any) bool {
		if p, ok := value.(*peer); ok {
			ch <- p
			return true
		}

		return false
	})

	return ch
}

// Query return connection interface based on socket parameter.
func (r *router) Query(id ID) *peer {
	// exist socket related peer?
	p, exists := r.table.Load(id)
	peer, ok := p.(*peer)

	if !exists || !ok {
		return nil
	}

	return peer
}

// Add forward method to internal sync.Map store for peer.
func (r *router) Add(peer *peer) {
	r.counter++
	r.table.Store(peer.ID(), peer)
}

// Len return the number of routed connections.
func (r *router) Len() uint8 {
	return r.counter
}

// Remove forward method to internal sync.Map to delete a connection from router.
// It return recently removed peer.
func (r *router) Remove(peer *peer) {
	r.counter--
	r.table.Delete(peer.ID())
}
