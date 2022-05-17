package network

// Message hold the information needed to exchange messages via pubsub events.
type Message struct {
	Type    Event      // Type of event published
	Payload []byte     // Custom data message published
	Peer    PeerStream // Peer interface sender
}

// Message factory
func NewMessage(event Event, payload []byte, peer PeerStream) *Message {
	return &Message{
		Peer:    peer,
		Type:    event,
		Payload: payload,
	}
}

// Reply message to sender peer
func (m *Message) Reply(msg []byte) (int, error) {
	return m.Peer.Send([]byte("Ping"))
}
