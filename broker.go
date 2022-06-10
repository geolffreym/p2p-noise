// Pubsub event notifications.
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

// Broker hash map event subscribers
type Broker struct {
	sync.RWMutex        // guards
	topics       Topics // topic subscriptions
}

func newBroker() *Broker {
	return &Broker{topics: make(Topics)}
}

// func (events *Broker) Unregister() error {
// 	if events.topics == nil {

// 	}
// }

// Register associate subscriber to a event channel;
func (events *Broker) Register(e Event, s *Subscriber) {
	// Mutex for writing topics.
	// Do not read while topics are written.
	// A blocked Lock call excludes new readers from acquiring the lock.
	// ref: https://pkg.go.dev/sync#RWMutex.Lock
	events.RWMutex.Lock()
	defer events.RWMutex.Unlock()
	events.topics.Add(e, s)
}

// Publish Emit/send concurrently messages to subscribers
func (events *Broker) Publish(msg Message) {
	// Mutex for reading topics.
	// Do not write while topics are read.
	// Write Lock can’t be acquired until all Read Locks are released.
	events.RWMutex.RLock()
	defer events.RWMutex.RUnlock()

	if _, ok := events.topics[msg.Type()]; ok {
		for _, subscriber := range events.topics[msg.Type()] {
			go func(s *Subscriber) {
				s.Emit(msg)
			}(subscriber)
		}
	}
}
