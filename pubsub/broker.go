package pubsub

type Event int
type Observer func(*Message) bool

const (
	SELF_LISTENING = iota
	NEWPEER_DETECTED
	MESSAGE_RECEIVED
)

type Channel map[Event][]*Subscriber

func (events Channel) Register(e Event, s *Subscriber) {
	// If not topic registered
	if _, ok := events[e]; !ok {
		events[e] = []*Subscriber{}
	}

	// Flag subscriber as subscribed
	s.topics[e] = true
	events[e] = append(events[e], s)
}

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
