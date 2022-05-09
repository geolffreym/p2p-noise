package pubsub

type Event int
type Observer func(*Message) bool

const (
	// Event for loopback on start listening event
	SELF_LISTENING = iota
	// Event to notify when a new peer connects
	NEWPEER_DETECTED
	// On new message received event
	MESSAGE_RECEIVED
)

// Hash map event subscribers
type Channel map[Event][]*Subscriber

// Associate subscriber to a Channel
// If channel event doesn't exist then is created
func (events Channel) Register(e Event, s *Subscriber) {
	// If not topic registered
	if _, ok := events[e]; !ok {
		events[e] = []*Subscriber{}
	}

	// Flag subscriber as subscribed
	s.topics[e] = true
	events[e] = append(events[e], s)
}

// Emit/send concurrently messages to subscribers
func (events Channel) Publish(msg *Message) {
	if _, ok := events[msg.Type]; ok {
		for _, subscriber := range events[msg.Type] {
			go func(s *Subscriber) {
				s.Emit(msg)
			}(subscriber)
		}
	}
}

// func (events Events) RemoveListener(e Event) {
// 	delete(events, e)
// }
