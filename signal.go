package noise

// signal hold the information needed to exchange messages via pubsub events.
type signal struct {
	event   Event  // Type of event published
	payload []byte // Custom data message published
}

// Type return Message event type published.
func (m signal) Type() Event { return m.event }

// Payload return custom data published.
func (m signal) Payload() []byte { return m.payload }

// SignalCtx keep message exchange context between network events.
// Each SignalCtx keep a state holding original Signal and related Peer.
type SignalCtx struct {
	signal signal
	peer   PeerCtx
}

// Payload forward internal signal event message payload.
func (s *SignalCtx) Payload() []byte {
	return s.signal.Payload()
}

// Type forward internal signal event message type.
func (s *SignalCtx) Type() Event {
	return s.signal.Type()
}

// Reply send an answer to contextual peer.
func (s *SignalCtx) Reply(msg []byte) (int, error) {
	return s.peer.Send(msg)
}
