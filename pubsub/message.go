package pubsub

type Message struct {
	Type    Event
	Payload []byte
}

func NewMessage(e Event, payload []byte) *Message {
	return &Message{
		Type:    e,
		Payload: payload,
	}
}
