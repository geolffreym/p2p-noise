package network

type Topic map[Event][]Subscriber

// Add append a new subscriber to event
// If topic event doesn't exist then is created.
func (t Topic) Add(e Event, s Subscriber) {
	// If not topic registered
	if _, ok := t[e]; !ok {
		t[e] = []Subscriber{}
	}

	t[e] = append(t[e], s)
}
