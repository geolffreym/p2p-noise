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

// SignalContext keep message exchange context between network events.
// Each SignalContext keep a state holding original Signal and related Peer.
type SignalContext struct {
	signal signal
	peer   *Peer
}

func newSignalContext(event Event, payload []byte, peer *Peer) SignalContext {
	return SignalContext{
		signal{event, payload},
		peer,
	}
}

// Payload forward internal signal event message payload.
func (s *SignalContext) Payload() []byte {
	return s.signal.Payload()
}

// Type forward internal signal event message type.
func (s *SignalContext) Type() Event {
	return s.signal.Type()
}

// Reply send an answer to contextual peer.
func (s *SignalContext) Reply(msg []byte) {
	s.peer.Send(msg)
}
