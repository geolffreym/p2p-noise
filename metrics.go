package noise

// metrics hold the statistics related to remote peers.
// We can add any method related to adaptive lookup logic here.
type metrics struct {
	handshakeStart uint32
	handshakeEnd   uint32
	latency        uint16
	bandwidth      uint16
}

// calculate weight
