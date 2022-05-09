package network

// Alias for string
type Socket string

type Router map[Socket]*Peer

//
func (r Router) Add(socket Socket, route *Peer) {
	r[socket] = route
}

func (r Router) Len() int {
	return len(r)
}

func (r Router) Delete(socket Socket) {
	delete(r, socket)
}
