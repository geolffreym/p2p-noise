package network

// Aliases to handle idiomatic `Socket` type
type Socket string

// Router hash table to associate Socket with Peers.
// eg. {127.0.0.1:4000: Peer}
type Router map[Socket]*Peer

// Add new socket => connection association
func (r Router) Add(socket Socket, peer *Peer) {
	r[socket] = peer
}

// Return the number of connections
func (r Router) Len() int {
	return len(r)
}

// Remove a connection from router
func (r Router) Delete(socket Socket) {
	delete(r, socket)
}
