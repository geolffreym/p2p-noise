package noise

// header keep the context for triggered signal.
type header struct {
	// Type of event published
	peer  *peer // Hold the involved peer
	event Event // Hold the triggered event
}

// Type return Event type published.
func (m header) Type() Event { return m.event }

// [Signal] it is a message interface to transport network events.
// Each Signal keep a immutable state holding original header and body.
type Signal struct {
	header header
	body   string
}

// Payload return an immutable payload.
func (s *Signal) Payload() string {
	return s.body
}

// Type forward internal signal header event type.
func (s *Signal) Type() Event {
	return s.header.Type()
}

// Reply send an answer to peer in context.
func (s *Signal) Reply(msg []byte) (uint, error) {
	return s.header.peer.Send(msg)
}
