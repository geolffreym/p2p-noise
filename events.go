package noise

import "context"

// Aliases to handle idiomatic `Event` type
type Event int

const (
	// Event to notify when a new peer get connected
	NewPeerDetected Event = iota
	// On new message received event
	MessageReceived
	// Closed peer connection
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
	body := body{peer.Socket().Bytes()}
	header := header{NewPeerDetected}
	signal := signal{header, body, peer}
	e.broker.Publish(signal)
}

// PeerDisconnected dispatch event peer disconnected.
func (e *events) PeerDisconnected(peer PeerCtx) {
	// Emit new notification
	body := body{peer.Socket().Bytes()}
	header := header{PeerDisconnected}
	signal := signal{header, body, peer}
	e.broker.Publish(signal)
}

// NewMessage dispatch event new message.
func (e *events) NewMessage(peer PeerCtx, msg []byte) {
	// Emit new notification
	body := body{msg}
	header := header{MessageReceived}
	signal := signal{header, body, peer}
	e.broker.Publish(signal)
}
