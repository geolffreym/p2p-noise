package network

// Aliases to handle idiomatic `string` type
type Socket string

// Router hash table to associate Socket with connection interfaces.
// eg.
// {127.0.0.1:4000: net.Conn}
type Router map[Socket]*Peer

// Add new socket => connection assoc
func (r Router) Add(socket Socket, route *Peer) {
	r[socket] = route
}

// Return the number of connections associated
func (r Router) Len() int {
	return len(r)
}

// Remove a connection from router
func (r Router) Delete(socket Socket) {
	delete(r, socket)
}
