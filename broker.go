package noise

import (
	"sync"
)

// IndexOf find index for element in slice.
// It return index if found else -1.
func IndexOf[T comparable](collection []T, el T) int {
	for i, v := range collection {
		if v == el {
			return i
		}
	}

	return -1
}

// topics `keep` registered events
type topics map[Event][]*subscriber

// Add append a new subscriber to event
// If topic event doesn't exist then is created.
func (t topics) Add(e Event, s *subscriber) {
	// If not topic registered
	if _, ok := t[e]; !ok {
		t[e] = []*subscriber{}
	}

	t[e] = append(t[e], s)
}

// Remove subscriber from topics
// It return true for removed subscriber from event else false.
func (t topics) Remove(e Event, s *subscriber) bool {
	// If not topic registered
	if _, ok := t[e]; ok {
		i := IndexOf(t[e], s)
		// if not match index for input subscriber
		if ^i == 0 {
			return false
		}

		t[e] = append(t[e][:i], t[e][i+1:]...)
		return true
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
func (b *broker) Register(e Event, s *subscriber) {
	// Lock while writing operation
	// If the lock is already in use, the calling goroutine blocks until the mutex is available.
	b.Mutex.Lock()
	b.topics.Add(e, s)
	b.Mutex.Unlock()
}

// Unregister remove associated subscriber from topics;
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
func (b *broker) Publish(msg SignalContext) uint8 {
	// Lock while reading operation
	b.Mutex.Lock()
	defer b.Mutex.Unlock()

	if _, ok := b.topics[msg.Type()]; ok {
		for _, sub := range b.topics[msg.Type()] {
			go func(s *subscriber) {
				s.Emit(msg)
			}(sub)
		}
		// Number of subscribers notified
		return uint8(len(b.topics[msg.Type()]))
	}

	return 0
}
