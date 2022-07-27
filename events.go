package noise

import "context"

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

type Broker interface {
	Register(e Event, s Subscriber)
	Publish(msg SignalCtx) uint8
}

type Subscriber interface {
	Emit(msg SignalCtx)
	Listen(ctx context.Context, ch chan<- SignalCtx)
}

type SignalCtx interface {
	Type() Event
	Payload() []byte
	Reply(msg []byte) (int, error)
}

// PeerCtx represents Peer in signal context.
// Each Signal keep a context with the peer involved in triggered event.
// eg. Signal{NewPeerDetected, PeerCtx}
type PeerCtx interface {
	Send(msg []byte) (int, error)
	Socket() Socket
}

type events struct {
	broker     Broker
	subscriber Subscriber
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
func (e *events) Subscriber() Subscriber {
	return e.subscriber
}

// PeerConnected dispatch event new peer detected.
func (e *events) PeerConnected(peer PeerCtx) {
	// Emit new notification
	addr := peer.Socket().Bytes()
	signal := signal{NewPeerDetected, addr, peer}
	e.broker.Publish(signal)
}

// PeerDisconnected dispatch event peer disconnected.
func (e *events) PeerDisconnected(peer PeerCtx) {
	// Emit new notification
	addr := peer.Socket().Bytes()
	signal := signal{PeerDisconnected, addr, peer}
	e.broker.Publish(signal)
}

// NewMessage dispatch event new message.
func (e *events) NewMessage(peer PeerCtx, msg []byte) {
	// Emit new notification
	signal := signal{MessageReceived, msg, peer}
	e.broker.Publish(signal)
}
