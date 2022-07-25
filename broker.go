package noise

import (
	"sync"
)

type data struct {
	s    []Subscriber
	sMap map[Subscriber]int
	size uint8
}

func (s *data) Len() uint8                { return s.size }
func (s *data) Subscribers() []Subscriber { return s.s }

// topics `keep` registered events
type topics map[Event]*data

// Add append a new subscriber to event
// If topic event doesn't exist then is created.
func (t topics) Add(e Event, s Subscriber) {
	// If not topic registered
	if _, ok := t[e]; !ok {
		t[e] = new(data)
		t[e].size = 0
		t[e].s = []Subscriber{}
		t[e].sMap = make(map[Subscriber]int)
	}

	t[e].s = append(t[e].s, s)
	t[e].sMap[s] = len(t[e].s) - 1
	t[e].size++
}

// Remove subscriber from topics
// It return true for removed subscriber from event else false.
func (t topics) Remove(e Event, s Subscriber) bool {
	// If topic registered
	if _, ok := t[e]; ok {
		if i, ok := t[e].sMap[s]; ok {
			// Clear topic from slice and map
			t[e].s = append(t[e].s[:i], t[e].s[i+1:]...)
			delete(t[e].sMap, s)
			t[e].size--
			return true
		}
	}

	return false
}

// broker hash map event subscribers
type broker struct {
	sync.Mutex        // guards
	topics     topics // topic subscriptions
}

func newBroker() *broker {
	return &broker{topics: make(topics)}
}

// Register associate subscriber to broker topics.
// It return new registered subscriber.
func (b *broker) Register(e Event, s Subscriber) {
	// Lock while writing operation
	// If the lock is already in use, the calling goroutine blocks until the mutex is available.
	b.Mutex.Lock()
	b.topics.Add(e, s)
	b.Mutex.Unlock()
}

// Unregister remove associated subscriber from topics;
// It return true for success else false.
func (b *broker) Unregister(e Event, s Subscriber) bool {
	// Lock while writing operation
	// If the lock is already in use, the calling goroutine blocks until the mutex is available.
	b.Mutex.Lock()
	defer b.Mutex.Unlock()       // This will be executed at the end of the enclosing function.
	return b.topics.Remove(e, s) // call first then return until then the mutex is available.
}

// Publish Emit/send concurrently messages to topic subscribers
// It return number of subscribers notified.
func (b *broker) Publish(msg SignalCtx) uint8 {
	// Lock while reading operation
	b.Mutex.Lock()
	defer b.Mutex.Unlock()

	if _, ok := b.topics[msg.Type()]; ok {
		topicData := b.topics[msg.Type()]
		topicLen := b.topics[msg.Type()].Len()
		subscribers := topicData.Subscribers()

		for _, sub := range subscribers {
			go func(s Subscriber) {
				s.Emit(msg)
			}(sub)
		}
		// Number of subscribers notified
		return topicLen
	}

	return 0
}
