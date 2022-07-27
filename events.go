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

// PeerCtx represents Peer in signal context.
// Each Signal keep a context with the peer involved in triggered event.
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

// Listen forward to Listen method to internal subscriber.
func (e *events) Listen(ctx context.Context, ch chan<- SignalCtx) {
	e.subscriber.Listen(ctx, ch)
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
