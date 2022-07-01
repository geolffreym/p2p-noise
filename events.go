package noise

// Aliases to handle idiomatic `Event` type
type Event int

const (
	// Event for loopback on start listening event
	SelfListening Event = iota
	// Event to notify when a new peer connects
	NewPeerDetected
	// On new message received event
	MessageReceived
	// On closed network
	ClosedConnection
	// Closed peer event
	PeerDisconnected
)

type events struct {
	broker     *broker
	subscriber *subscriber
}

func newEvents() *events {
	subscriber := newSubscriber()
	broker := newBroker()
	// register default events
	broker.Register(SelfListening, subscriber)
	broker.Register(NewPeerDetected, subscriber)
	broker.Register(MessageReceived, subscriber)
	broker.Register(ClosedConnection, subscriber)
	broker.Register(PeerDisconnected, subscriber)

	return &events{
		broker,
		subscriber,
	}
}

// Subscriber return event subscriber interface.
func (e *events) Subscriber() *subscriber {
	return e.subscriber
}

// PeerConnected dispatch event new peer detected.
func (e *events) PeerConnected(addr []byte) {
	// Emit new notification
	e.broker.Publish(Message{
		NewPeerDetected,
		addr,
	})
}

// PeerDisconnected dispatch event peer disconnected.
func (e *events) PeerDisconnected(addr []byte) {
	// Emit new notification
	e.broker.Publish(Message{
		PeerDisconnected,
		addr,
	})
}

// Listening dispatch event self listening.
func (e *events) Listening(addr []byte) {
	// Emit new notification
	e.broker.Publish(Message{
		SelfListening,
		addr,
	})
}

// NewMessage dispatch event new message.
func (e *events) NewMessage(msg []byte) {
	// Emit new notification
	e.broker.Publish(Message{
		MessageReceived,
		msg,
	})
}

// ClosedConnection dispatch event closed connection.
func (e *events) ClosedConnection() {
	// Emit new notification
	e.broker.Publish(Message{
		ClosedConnection,
		[]byte(""),
	})
}
