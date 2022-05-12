package network

// Message hold the information needed to exchange messages via pubsub events.
type Message struct {
	Type    Event  // Type of event published
	From    *Peer  // Peer interface sender
	Payload []byte // Custom data message published
}

// Message factory
func NewMessage(event Event, payload []byte, from *Peer) *Message {
	return &Message{
		From:    from,
		Type:    event,
		Payload: payload,
	}
}
