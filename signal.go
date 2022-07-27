package noise

// signal keep message exchange context between network events.
// Each signal keep a state holding original Signal and related Peer.
type signal struct {
	// Type of event published
	event Event
	// Custom data message published
	payload []byte
	// Related peer interface
	peer PeerCtx
}

// Type return Message event type published.
func (m signal) Type() Event { return m.event }

// Payload return custom data published.
func (m signal) Payload() []byte { return m.payload }

// Reply send an answer to peer in context.
func (s signal) Reply(msg []byte) (int, error) {
	return s.peer.Send(msg)
}
