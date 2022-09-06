package noise

import "unsafe"

// header keep the context for triggered signal.
type header struct {
	// Type of event published
	peer  *peer // Hold the involved peer
	event Event // Hold the triggered event
}

// Type return Event type published.
func (m header) Type() Event { return m.event }

// body keep payload related to signal.
type body []byte

// Payload return custom data published.
func (m body) Bytes() []byte { return m }
func (m body) String() string {
	// A pointer value can't be converted to an arbitrary pointer type.
	// ref: https://go101.org/article/unsafe.html
	// no-copy conversion
	// ref: https://github.com/golang/go/issues/25484
	return *(*string)(unsafe.Pointer(&m))
}

// [Signal] it is a message interface to transport network events.
// Each Signal keep a immutable state holding original header and body.
type Signal struct {
	header header
	body   body
}

// Payload forward internal signal body payload.
// Return an immutable string payload.
func (s *Signal) Payload() string {
	return s.body.String()
}

// Type forward internal signal header event type.
func (s *Signal) Type() Event {
	return s.header.Type()
}

// Reply send an answer to peer in context.
func (s *Signal) Reply(msg []byte) (int, error) {
	return s.header.peer.Send(msg)
}
