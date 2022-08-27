package noise

import "context"

// [Event] aliases for int type.
type Event int

const (
	// Event to notify when a new peer get connected
	NewPeerDetected Event = iota
	// On new message received event
	MessageReceived
	// Closed peer connection
	PeerDisconnected
)

// events handle event exchange between [Node] and network.
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
func (e *events) Listen(ctx context.Context, ch chan<- Signal) {
	e.subscriber.Listen(ctx, ch)
}

// PeerConnected dispatch event when new peer is detected.
func (e *events) PeerConnected(peer *peer) {
	// Emit new notification
	body := body{peer.Socket().Bytes()}
	header := header{NewPeerDetected}
	signal := Signal{header, body, peer}
	e.broker.Publish(signal)
}

// PeerDisconnected dispatch event when peer get disconnected.
func (e *events) PeerDisconnected(peer *peer) {
	// Emit new notification
	body := body{peer.Socket().Bytes()}
	header := header{PeerDisconnected}
	signal := Signal{header, body, peer}
	e.broker.Publish(signal)
}

// NewMessage dispatch event when a new message is received.
func (e *events) NewMessage(peer *peer, msg []byte) {
	// Emit new notification
	body := body{msg}
	header := header{MessageReceived}
	signal := Signal{header, body, peer}
	e.broker.Publish(signal)
}
