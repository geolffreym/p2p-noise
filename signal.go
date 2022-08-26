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

// signal implements Signal interface.
// Each signal keep a state holding original header, body and related peer.
type signal struct {
	header header
	body   body
	peer   *peer
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
