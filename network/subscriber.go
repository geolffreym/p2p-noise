package network

import (
	"sync"
)

// Subscribed map event list for subscriber
type Subscriptions map[Event]bool

// Set event as subscribed = true
func (s Subscriptions) Add(event Event) {
	s[event] = true
}

// Subscriber synchronize messages from events
// and keep a record of current subscriptions for events.
type Subscriber struct {
	mutex      sync.RWMutex  // Mutual exclusion
	message    chan *Message // Message exchange channel
	subscribed Subscriptions // Keep tracking of subscribed events
}

// Subscriber factory
func NewSubscriber() *Subscriber {
	return &Subscriber{
		message:    make(chan *Message),
		subscribed: make(Subscriptions),
	}
}

// Send Message to channel buffer.
func (s *Subscriber) Emit(msg *Message) {
	// Get lock to enforce sync order messages
	// No read messages while writing
	// https://gobyexample.com/mutexes
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	s.message <- msg
}

// Listen and wait fot Message synchronization from channel.
// When a new message is added to channel buffer the
// Observer is executed with new Message propagated as param.
// !Important: If Observer returns false the routine stop "listening".
func (s *Subscriber) Listen(cb Observer) {
	go func(subscriber *Subscriber, call Observer) {
		for {
			if msg, ok := <-subscriber.message; ok {
				keepAlive := call(msg)
				// Close if callback returns false.
				if keepAlive == false {
					return
				}
			}
		}
	}(s, cb)
}
