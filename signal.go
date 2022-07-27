package noise

// header keep message exchange context between network events.
// Each header keep a state holding original Signal and related Peer.
type header struct {
	// Type of event published
	event Event
}

// Type return Message event type published.
func (m header) Type() Event { return m.event }

type body struct {
	// Custom data message published
	payload []byte
}

// Payload return custom data published.
func (m body) Payload() []byte { return m.payload }

// signal keep message exchange context between network events.
// Each signal keep a state holding original header, body and related peer.
// Anatomy of a signal.
type signal struct {
	header header
	body   body
	peer   PeerCtx
}

// Payload forward internal signal event message payload.
func (s signal) Payload() []byte {
	return s.body.Payload()
}

// Type forward internal signal event message type.
func (s signal) Type() Event {
	return s.header.Type()
}

// Reply send an answer to contextual peer.
func (s signal) Reply(msg []byte) (int, error) {
	return s.peer.Send(msg)
}
