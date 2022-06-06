// Pubsub event notifications.
package network

import "sync"

// Aliases to handle idiomatic `Event` type
type Event int

// Callback interface to propagate events notifications
type Observer func(Message) bool

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

type Events interface {
	Register(e Event, s Subscriber)
	Publish(msg Message)
	Topics() Topic
}

// Events hash map event subscribers
type events struct {
	sync.RWMutex       // guards
	topics       Topic // topic subscriptions
}

func NewEvents() Events {
	return &events{topics: make(Topic, 0)}
}

func (events *events) Topics() Topic { return events.topics }

// Register associate subscriber to a event channel;
func (events *events) Register(e Event, s Subscriber) {
	// Mutex for writing topics.
	// Do not read while topics are written.
	// A blocked Lock call excludes new readers from acquiring the lock.
	// ref: https://pkg.go.dev/sync#RWMutex.Lock
	events.RWMutex.Lock()
	defer events.RWMutex.Unlock()
	events.topics.Add(e, s)
}

// Publish Emit/send concurrently messages to subscribers
func (events *events) Publish(msg Message) {
	// Mutex for reading topics.
	// Do not write while topics are read.
	// Write Lock can’t be acquired until all Read Locks are released.
	events.RWMutex.RLock()
	defer events.RWMutex.RUnlock()

	if _, ok := events.topics[msg.Type()]; ok {
		for _, subscriber := range events.topics[msg.Type()] {
			go func(s Subscriber) {
				s.Emit(msg)
			}(subscriber)
		}
	}
}
