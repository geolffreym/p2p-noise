// Pubsub event notifications.
package network

import "sync"

// Aliases to handle idiomatic `Event` type
type Event int

// Callback interface to propagate events notifications
type Observer func(*Message) bool

const (
	// Event for loopback on start listening event
	SELF_LISTENING = iota
	// Event to notify when a new peer connects
	NEWPEER_DETECTED
	// On new message received event
	MESSAGE_RECEIVED
	// On closed network
	CLOSED_CONNECTION
)

type Topic map[Event][]*Subscriber

// Add append a new subscriber to event
// If topic event doesn't exist then is created.
func (t Topic) Add(e Event, s *Subscriber) {
	// If not topic registered
	if _, ok := t[e]; !ok {
		t[e] = []*Subscriber{}
	}

	t[e] = append(t[e], s)
}

// Events hash map event subscribers
type Events struct {
	sync.RWMutex       // guards
	topics       Topic // topic subscriptions
}

func NewEvents() *Events {
	return &Events{topics: make(Topic)}
}

// Register associate subscriber to a event channel;
func (events *Events) Register(e Event, s *Subscriber) {
	// Mutex for writing topics.
	// Do not read while topics are written.
	// If a goroutine holds a RWMutex for reading and another goroutine might call Lock,
	// no goroutine should expect to be able to acquire a read lock until the initial read lock is released.
	// ref: https://pkg.go.dev/sync#RWMutex.Lock
	events.RWMutex.Lock()
	defer events.RWMutex.Unlock()
	events.topics.Add(e, s)
}

// Publish Emit/send concurrently messages to subscribers
func (events *Events) Publish(msg *Message) {
	// Mutex for reading topics.
	// Do not write while topics are read.
	// Write Lock can’t be acquired until all Read Locks are released.
	events.RWMutex.RLock()
	defer events.RWMutex.RUnlock()

	if _, ok := events.topics[msg.Type]; ok {
		for _, subscriber := range events.topics[msg.Type] {
			go func(s *Subscriber) {
				s.Emit(msg)
			}(subscriber)
		}
	}
}
