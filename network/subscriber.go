package network

import (
	"sync"
)

type SubscriberListener interface {
	Listen(cb Observer)
}

type SubscriberEmitter interface {
	Emit(msg Message)
}

type SubscriberMessenger interface {
	Message() chan Message
}

type Subscriber interface {
	SubscriberEmitter
	SubscriberListener
	SubscriberMessenger
}

// Subscriber synchronize messages from events
// and keep a record of current subscriptions for events.
type subscriber struct {
	sync.RWMutex              // Mutual exclusion
	message      chan Message // Message exchange channel
}

// Subscriber factory
func NewSubscriber() Subscriber {
	return &subscriber{
		message: make(chan Message),
	}
}

func (s *subscriber) Message() chan Message {
	return s.message
}

// Send Message to channel buffer.
func (s *subscriber) Emit(msg Message) {
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
func (s *subscriber) Listen(cb Observer) {
	go func(subscriber *subscriber, call Observer) {
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
