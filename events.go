package noise

import (
	"context"
	"unsafe"
)

// byteToString convert an array of bytes to a string with no-copy strategy.
func bytesToString(b []byte) string {
	// Optimizing space with ordered types.
	// perf: no allocation/copy to convert to string instead take the already existing byte slice to create a string struct.
	// WARNING: use this approach with caution and only if we are sure that the bytes slice is not gonna change.
	return *(*string)(unsafe.Pointer(&b))
}

// [Event] aliases for int type.
type Event uint8

const (
	// Event to notify when a new peer get connected
	NewPeerDetected Event = iota
	// On new message received event
	MessageReceived
	// Closed peer connection
	PeerDisconnected
	// Emitted when the node is ready to accept incoming connections
	SelfListening
)

// events handle event exchange between [Node] and network.
type events struct {
	broker     *broker
	subscriber *subscriber
}

func newEvents() *events {
	subscriber := newSubscriber()
	// !IMPORTANT if new events are added the size should be equal to new events number.
	// we need only 4 spaces one for each event, adding this avoids potential map growth.
	// https://100go.co/#inefficient-map-initialization-27
	broker := newBroker(4)
	// register default events
	broker.Register(NewPeerDetected, subscriber)
	broker.Register(MessageReceived, subscriber)
	broker.Register(PeerDisconnected, subscriber)
	broker.Register(SelfListening, subscriber)

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
	body := peer.ID().String()
	header := header{peer, NewPeerDetected}
	signal := Signal{header, body}
	e.broker.Publish(signal)
}

// PeerDisconnected dispatch event when peer get disconnected.
func (e *events) PeerDisconnected(peer *peer) {
	// Emit new notification
	body := peer.ID().String()
	header := header{peer, PeerDisconnected}
	signal := Signal{header, body}
	e.broker.Publish(signal)
}

// SelfListening dispatch event when node is ready.
func (e *events) SelfListening(addr string) {
	// Emit new notification
	header := header{nil, SelfListening}
	signal := Signal{header, addr}
	e.broker.Publish(signal)
}

// NewMessage dispatch event when a new message is received.
func (e *events) NewMessage(peer *peer, msg []byte) {
	// Emit new notification
	message := bytesToString(msg)
	header := header{peer, MessageReceived}
	signal := Signal{header, message}
	e.broker.Publish(signal)
}
