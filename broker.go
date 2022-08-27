package noise

import (
	"sync"
)

type topic struct {
	s    []*subscriber
	sMap map[*subscriber]int
	size uint8
}

func (s *topic) Len() uint8                 { return s.size }
func (s *topic) Subscribers() []*subscriber { return s.s }

// topics `keep` registered events
type topics map[Event]*topic

// Add append a new subscriber to event
// If topic event doesn't exist then is created.
func (t topics) Add(e Event, s *subscriber) {
	// If not topic registered
	if _, ok := t[e]; !ok {
		t[e] = new(topic)
		t[e].size = 0
		t[e].s = []*subscriber{}
		t[e].sMap = make(map[*subscriber]int)
	}

	t[e].s = append(t[e].s, s)
	t[e].sMap[s] = len(t[e].s) - 1
	t[e].size++
}

// Remove subscriber from topics
// It return true for removed subscriber from event else false.
func (t topics) Remove(e Event, s *subscriber) bool {
	// Is topic registered?
	to, topicExists := t[e]
	if !topicExists {
		return false
	}

	// Check if subscriber exists in topic
	i, isSubscribed := to.sMap[s]
	if !isSubscribed {
		return false
	}

	// Clear topic from slice and map
	to.s = append(to.s[:i], to.s[i+1:]...) // re-slice to remove old subscriber
	delete(to.sMap, s)                     // remove index from mapping
	to.size--                              // reduce size of subscribers in topic
	return true

}

// broker exchange messages between events and subscriber.
// Each broker receive published signal from event for later emit it to subscriber.
type broker struct {
	sync.Mutex        // guards
	topics     topics // topic subscriptions
}

func newBroker() *broker {
	return &broker{topics: make(topics)}
}

// Register associate subscriber to broker topics.
// It return new registered subscriber.
func (b *broker) Register(e Event, s *subscriber) {
	// Lock while writing operation
	// If the lock is already in use, the calling goroutine blocks until the mutex is available.
	b.Mutex.Lock()
	b.topics.Add(e, s)
	b.Mutex.Unlock()
}

// Unregister remove associated subscriber from topics.
// It return true for success else false.
func (b *broker) Unregister(e Event, s *subscriber) bool {
	// Lock while writing operation
	// If the lock is already in use, the calling goroutine blocks until the mutex is available.
	b.Mutex.Lock()
	defer b.Mutex.Unlock()       // This will be executed at the end of the enclosing function.
	return b.topics.Remove(e, s) // call first then return until then the mutex is available.
}

// Publish Emit/send concurrently messages to topic subscribers
// It return number of subscribers notified.
func (b *broker) Publish(msg Signal) uint8 {
	// Lock while reading operation
	b.Mutex.Lock()
	defer b.Mutex.Unlock()

	// Check if topic is registered before try to emit messages to subscribers.
	topicData, topicExists := b.topics[msg.Type()]
	if !topicExists {
		return 0
	}

	// How many subscribers exists in topic?
	topicLen := topicData.Len()
	// Subscribers in topic!!
	subscribers := topicData.Subscribers() // Subscribers in topic!!
	// For each subscriber in topic registered emit a new signal
	for _, sub := range subscribers {
		go func(s *subscriber) {
			s.Emit(msg)
		}(sub)
	}
	// Number of subscribers notified
	return topicLen

}
