package noise

// metrics hold the statistics related to remote peers.
// We can add any method related to adaptive metrics logic here.
type metrics struct {
	handshakeStart uint64
	handshakeEnd   uint64
	latency        uint16
	bandwidth      uint16
}

// calculate weight
