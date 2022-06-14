package noise

// Message hold the information needed to exchange messages via pubsub events.
type Message struct {
	payload []byte // Custom data message published
	event   Event  // Type of event published
}

// Message factory
func newMessage(event Event, payload []byte) Message {
	return Message{
		event:   event,
		payload: payload,
	}
}

// Type return Message event type published.
func (m Message) Type() Event { return m.event }

// Payload return custom data published.
func (m Message) Payload() []byte { return m.payload }
