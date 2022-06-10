package noise

import (
	"context"
)

// Aliases to handle idiomatic `Event` type
type (
	Event    int
	Observer func(Message)
)

const (
	// Event for loopback on start listening event
	SELF_LISTENING Event = iota
	// Event to notify when a new peer connects
	NEWPEER_DETECTED
	// On new message received event
	MESSAGE_RECEIVED
	// On closed network
	CLOSED_CONNECTION
	// Closed peer event
	PEER_DISCONNECTED
)

type Events struct {
	broker     *Broker
	subscriber *Subscriber
}

func newMessenger() *Events {
	broker := newBroker()
	subscriber := newSubscriber()

	broker.Register(SELF_LISTENING, subscriber)
	broker.Register(NEWPEER_DETECTED, subscriber)
	broker.Register(MESSAGE_RECEIVED, subscriber)
	broker.Register(CLOSED_CONNECTION, subscriber)
	broker.Register(PEER_DISCONNECTED, subscriber)

	return &Events{
		broker:     broker,
		subscriber: subscriber,
	}
}

func (msgr *Events) Listen(cb Observer) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	msgr.subscriber.Listen(ctx, cb)
	return cancel
}

// dispatch event new peer detected
func (msgr *Events) PeerConnected(addr []byte) {
	// Emit new notification
	message := newMessage(NEWPEER_DETECTED, addr)
	msgr.broker.Publish(message)
}

// dispatch event peer disconnected
func (msgr *Events) PeerDisconnected(addr []byte) {
	// Emit new notification
	message := newMessage(PEER_DISCONNECTED, addr)
	msgr.broker.Publish(message)
}

// dispatch event self listening
func (msgr *Events) Listening(addr []byte) {
	// Emit new notification
	message := newMessage(SELF_LISTENING, addr)
	msgr.broker.Publish(message)
}

// dispatch event new message
func (msgr *Events) NewMessage(msg []byte) {
	// Emit new notification
	message := newMessage(MESSAGE_RECEIVED, msg)
	msgr.broker.Publish(message)
}

// dispatch event closed connection
func (msgr *Events) ClosedConnection() {
	// Emit new notification
	message := newMessage(CLOSED_CONNECTION, []byte(""))
	msgr.broker.Publish(message)
}
