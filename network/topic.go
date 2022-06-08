package network

type Topics map[Event][]Messenger

func NewTopic() Topics {
	return make(Topics)
}

// Add append a new subscriber to event
// If topic event doesn't exist then is created.
func (t Topics) Add(e Event, s Messenger) {
	// If not topic registered
	if _, ok := t[e]; !ok {
		t[e] = []Messenger{}
	}

	t[e] = append(t[e], s)
}
