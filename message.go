package noise

// Message hold the information needed to exchange messages via pubsub events.
type Message struct {
	event   Event  // Type of event published
	payload []byte // Custom data message published
}

// Type return Message event type published.
func (m Message) Type() Event { return m.event }

// Payload return custom data published.
func (m Message) Payload() []byte { return m.payload }
