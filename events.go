package noise

// Aliases to handle idiomatic `Event` type
type Event int

const (
	// Event to notify when a new peer connects
	NewPeerDetected Event = iota
	// On new message received event
	MessageReceived
	// Closed peer event
	PeerDisconnected
)

type PeerCtx interface {
	Send(msg []byte) (int, error)
	Socket() Socket
}

type events struct {
	broker     *broker
	subscriber *subscriber
}

func newEvents() *events {
	subscriber := newSubscriber()
	broker := newBroker()
	// register default events
	broker.Register(NewPeerDetected, subscriber)
	broker.Register(MessageReceived, subscriber)
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
func (e *events) PeerConnected(peer PeerCtx) {
	// Emit new notification
	addr := peer.Socket().Bytes()
	signal := signal{NewPeerDetected, addr}
	e.broker.Publish(SignalCtx{signal, peer})
}

// PeerDisconnected dispatch event peer disconnected.
func (e *events) PeerDisconnected(peer PeerCtx) {
	// Emit new notification
	addr := peer.Socket().Bytes()
	signal := signal{PeerDisconnected, addr}
	e.broker.Publish(SignalCtx{signal, peer})
}

// NewMessage dispatch event new message.
func (e *events) NewMessage(peer PeerCtx, msg []byte) {
	// Emit new notification
	signal := signal{MessageReceived, msg}
	e.broker.Publish(SignalCtx{signal, peer})
}
