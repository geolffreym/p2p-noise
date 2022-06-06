package network

type Message interface {
	Reply(msg []byte) (int, error)
	Peer() PeerStreamer
	Payload() []byte
	Type() Event
}

// Message hold the information needed to exchange messages via pubsub events.
type message struct {
	event   Event        // Type of event published
	payload []byte       // Custom data message published
	peer    PeerStreamer // Peer interface sender
}

// Message factory
func NewMessage(event Event, payload []byte, peer PeerStreamer) Message {
	return &message{
		peer:    peer,
		event:   event,
		payload: payload,
	}
}

func (m *message) Peer() PeerStreamer { return m.peer }
func (m *message) Type() Event        { return m.event }
func (m *message) Payload() []byte    { return m.payload }

// Reply message to sender peer
func (m *message) Reply(msg []byte) (int, error) {
	return m.peer.Send(msg)
}
