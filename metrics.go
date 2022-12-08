package noise

// metrics hold the statistics related to remote peers.
// We can add any method related to adaptive lookup logic here.
// Please see [docs] for more information.
//
// [docs]: https://arxiv.org/pdf/1509.04417.pdf
type metrics struct {
	handshakeTime uint32 // how long took the handshake to complete.
	latency       uint16 // rtt in ms
	bandwidth     uint16 // remote peer bandwidth
	nonce         uint16 // nonce ordering factor
	sent          uint16 // sent messages
	recv          uint16 // received messages
}

// calculate weight
// builder pattern?
