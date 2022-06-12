package noise

import (
	"sync"
)

// Callback interface to propagate events notifications
type Topics map[Event][]*Subscriber

// Add append a new subscriber to event
// If topic event doesn't exist then is created.
func (t Topics) Add(e Event, s *Subscriber) {
	// If not topic registered
	if _, ok := t[e]; !ok {
		t[e] = []*Subscriber{}
	}

	t[e] = append(t[e], s)
}

// Remove subscriber from topics
// It return true for removed subscriber from event else false.
func (t Topics) Remove(e Event, s *Subscriber) bool {
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

// Broker hash map event subscribers
type Broker struct {
	sync.Mutex        // guards
	topics     Topics // topic subscriptions
}

func newBroker() *Broker {
	return &Broker{topics: make(Topics)}
}

// Register associate subscriber to broker topics.
// It return new registered subscriber.
func (b *Broker) Register(e Event, s *Subscriber) *Subscriber {
	// Lock while writing operation
	// If the lock is already in use, the calling goroutine blocks until the mutex is available.
	b.Mutex.Lock()
	defer b.Mutex.Unlock()
	b.topics.Add(e, s)
	return s
}

// Unregister remove associated subscriber from topics;
// It return true for success else false.
func (b *Broker) Unregister(e Event, s *Subscriber) bool {
	// Lock while writing operation
	// If the lock is already in use, the calling goroutine blocks until the mutex is available.
	b.Mutex.Lock()
	defer b.Mutex.Unlock()
	return b.topics.Remove(e, s)
}

// Publish Emit/send concurrently messages to topic subscribers
// It return number of subscribers notified.
func (b *Broker) Publish(msg Message) int {
	// Lock while reading operation
	b.Mutex.Lock()
	defer b.Mutex.Unlock()

	if _, ok := b.topics[msg.Type()]; ok {
		for _, subscriber := range b.topics[msg.Type()] {
			go func(s *Subscriber) {
				s.Emit(msg)
			}(subscriber)
		}
		// Number of subscribers notified
		return len(b.topics[msg.Type()])
	}

	return 0
}
