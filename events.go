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

type Events struct {
	broker     *Broker
	subscriber *Subscriber
}

func newEvents() *Events {
	subscriber := newSubscriber()
	broker := newBroker()
	// register default events
	broker.Register(SelfListening, subscriber)
	broker.Register(NewPeerDetected, subscriber)
	broker.Register(MessageReceived, subscriber)
	broker.Register(ClosedConnection, subscriber)
	broker.Register(PeerDisconnected, subscriber)

	return &Events{
		broker:     broker,
		subscriber: subscriber,
	}
}

// Subscriber return event subscriber interface.
func (e *Events) Subscriber() *Subscriber {
	return e.subscriber
}

// PeerConnected dispatch event new peer detected.
func (e *Events) PeerConnected(addr []byte) {
	// Emit new notification
	message := newMessage(NewPeerDetected, addr)
	e.broker.Publish(message)
}

// PeerDisconnected dispatch event peer disconnected.
func (e *Events) PeerDisconnected(addr []byte) {
	// Emit new notification
	message := newMessage(PeerDisconnected, addr)
	e.broker.Publish(message)
}

// Listening dispatch event self listening.
func (e *Events) Listening(addr []byte) {
	// Emit new notification
	message := newMessage(SelfListening, addr)
	e.broker.Publish(message)
}

// NewMessage dispatch event new message.
func (e *Events) NewMessage(msg []byte) {
	// Emit new notification
	message := newMessage(MessageReceived, msg)
	e.broker.Publish(message)
}

// ClosedConnection dispatch event closed connection.
func (e *Events) ClosedConnection() {
	// Emit new notification
	message := newMessage(ClosedConnection, []byte(""))
	e.broker.Publish(message)
}
