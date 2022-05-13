package network

// Message hold the information needed to exchange messages via pubsub events.
type Message struct {
	Type    Event  // Type of event published
	Peer    *Peer  // Peer interface sender
	Payload []byte // Custom data message published
}

// Message factory
func NewMessage(event Event, payload []byte, peer *Peer) *Message {
	return &Message{
		Peer:    peer,
		Type:    event,
		Payload: payload,
	}
}

// Reply message to sender peer
func (m *Message) Reply(msg []byte) {
	m.Peer.Send([]byte("Ping"))
}
