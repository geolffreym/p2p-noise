package noise

// Message hold the information needed to exchange messages via pubsub events.
type Message struct {
	payload []byte // Custom data message published
	ID      Event  // Type of event published
}

// Message factory
func newMessage(ID Event, payload []byte) Message {
	return Message{
		ID:      ID,
		payload: payload,
	}
}

func (m Message) Type() Event     { return m.ID }
func (m Message) Payload() []byte { return m.payload }
