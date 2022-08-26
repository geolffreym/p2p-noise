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

// [Signal] it is a message interface to transport network events.
type Signal interface {
	// Return event message type.
	Type() Event
	// Return event message payload.
	Payload() []byte
	// Reply send an answer to peer in context.
	Reply(msg []byte) (int, error)
}

// [Subscriber] intercept Signals from already subscribed topics in [Broker].
type Subscriber interface {
	// Emit signal using not-buffered channel.
	Emit(msg Signal)
	// Listen and wait for Signal synchronization from channel.
	Listen(ctx context.Context, ch chan<- Signal)
}

// [Broker] exchange messages between [Events] and [Subscriber].
// Each [Broker] receive published [Signal] from [Event] for later emit it to [Subscriber].
type Broker interface {
	// Register associate Subscriber to broker topics.
	Register(e Event, s Subscriber)
	// Unregister remove associated subscriber from topics.
	Unregister(e Event, s Subscriber) bool
	// Publish Emit/send concurrently messages to topic subscribers.
	Publish(msg Signal) uint8
}

// events implements Events interface.
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

// Listen forward to Listen method to internal subscriber.
func (e *events) Listen(ctx context.Context, ch chan<- Signal) {
	e.subscriber.Listen(ctx, ch)
}

// PeerConnected dispatch event when new peer is detected.
func (e *events) PeerConnected(peer Peer) {
	// Emit new notification
	body := body{peer.Socket().Bytes()}
	header := header{NewPeerDetected}
	signal := signal{header, body, peer}
	e.broker.Publish(signal)
}

// PeerDisconnected dispatch event when peer get disconnected.
func (e *events) PeerDisconnected(peer Peer) {
	// Emit new notification
	body := body{peer.Socket().Bytes()}
	header := header{PeerDisconnected}
	signal := signal{header, body, peer}
	e.broker.Publish(signal)
}

// NewMessage dispatch event when a new message is received.
func (e *events) NewMessage(peer Peer, msg []byte) {
	// Emit new notification
	body := body{msg}
	header := header{MessageReceived}
	signal := signal{header, body, peer}
	e.broker.Publish(signal)
}
