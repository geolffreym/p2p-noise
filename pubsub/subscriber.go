package pubsub

import "sync"

// Subscriber synchronize messages from events
// and keep a record of current subscriptions for events.
type Subscriber struct {
	mutex   sync.Mutex     // Mutual exclusion
	message chan *Message  // Message exchange channel
	events  map[Event]bool // Keep tracking of subscribed events
}

// Subscriber factory
func NewSubscriber() *Subscriber {
	return &Subscriber{
		message: make(chan *Message),
		events:  map[Event]bool{},
	}
}

// Send Message to channel buffer.
func (s *Subscriber) Emit(msg *Message) {
	// Get lock to enforce sync order messages
	// https://gobyexample.com/mutexes
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.message <- msg
}

// Listen and wait fot Message synchronization from channel.
// When a new message is added to channel buffer the
// Observer is executed with new Message propagated as param.
// !Important: If Observer returns false the routine stop "listening".
func (s *Subscriber) Listen(cb Observer) {
	go func(call Observer) {
		for {
			if msg, ok := <-s.message; ok {
				keepAlive := call(msg)
				// Close if callback returns false.
				if keepAlive == false {
					return
				}
			}
		}
	}(cb)
}
