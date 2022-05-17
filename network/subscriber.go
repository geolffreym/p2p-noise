package network

import (
	"sync"
)

// Subscriber synchronize messages from events
// and keep a record of current subscriptions for events.
type Subscriber struct {
	sync.RWMutex               // Mutual exclusion
	message      chan *Message // Message exchange channel
}

// Subscriber factory
func NewSubscriber() *Subscriber {
	return &Subscriber{
		message: make(chan *Message),
	}
}

// Send Message to channel buffer.
func (s *Subscriber) Emit(msg *Message) {
	// Lock exclusive writing
	// https://gobyexample.com/mutexes
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()
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
