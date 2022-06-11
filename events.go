package noise

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

func newEvents() *Events {
	subscriber := newSubscriber()
	broker := newBroker()
	// register default events
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

// Getter for event subscriber interface
func (events *Events) Subscriber() *Subscriber {
	return events.subscriber
}

// dispatch event new peer detected
func (events *Events) PeerConnected(addr []byte) {
	// Emit new notification
	message := newMessage(NEWPEER_DETECTED, addr)
	events.broker.Publish(message)
}

// dispatch event peer disconnected
func (events *Events) PeerDisconnected(addr []byte) {
	// Emit new notification
	message := newMessage(PEER_DISCONNECTED, addr)
	events.broker.Publish(message)
}

// dispatch event self listening
func (events *Events) Listening(addr []byte) {
	// Emit new notification
	message := newMessage(SELF_LISTENING, addr)
	events.broker.Publish(message)
}

// dispatch event new message
func (events *Events) NewMessage(msg []byte) {
	// Emit new notification
	message := newMessage(MESSAGE_RECEIVED, msg)
	events.broker.Publish(message)
}

// dispatch event closed connection
func (events *Events) ClosedConnection() {
	// Emit new notification
	message := newMessage(CLOSED_CONNECTION, []byte(""))
	events.broker.Publish(message)
}
