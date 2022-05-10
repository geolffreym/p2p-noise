package pubsub

// Message hold the information needed to exchange messages via pubsub events.
type Message struct {
	Type    Event  // Type of event published
	Payload []byte // Custom data message published
}

// Message factory
func NewMessage(e Event, payload []byte) *Message {
	return &Message{
		Type:    e,
		Payload: payload,
	}
}
