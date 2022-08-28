package noise

// header keep event type related to signal.
type header struct {
	// Type of event published
	event Event
}

// Type return Event type published.
func (m header) Type() Event { return m.event }

// body keep payload related to signal.
type body struct {
	// Custom data message published
	payload []byte
}

// Payload return custom data published.
func (m body) Payload() []byte { return m.payload }

// [Signal] it is a message interface to transport network events.
// Each Signal keep a state holding original header, body and related peer.
type Signal struct {
	header header
	body   body
	// Use a pointer if you are using a type that has methods with pointer receivers.
	peer *peer
}

// Payload forward internal signal body payload.
func (s *Signal) Payload() []byte {
	return s.body.Payload()
}

// Type forward internal signal header event type.
func (s *Signal) Type() Event {
	return s.header.Type()
}

// Reply send an answer to peer in context.
func (s *Signal) Reply(msg []byte) (int, error) {
	return s.peer.Send(msg)
}
