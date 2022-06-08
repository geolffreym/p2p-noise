package network

import (
	"context"
	"sync"
)

type Messenger interface {
	Listen(ctx context.Context, cb Observer)
	Message() chan Message
	Emit(msg Message)
}

// Subscriber synchronize messages from events
// and keep a record of current subscriptions for events.
type messenger struct {
	sync.RWMutex              // Mutual exclusion
	message      chan Message // Message exchange channel
}

// Subscriber factory
func NewMessenger() Messenger {
	return &messenger{
		message: make(chan Message),
	}
}

func (s *messenger) Message() chan Message {
	return s.message
}

// Send Message to channel buffer.
func (s *messenger) Emit(msg Message) {
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
func (s *messenger) Listen(ctx context.Context, cb Observer) {
	go func(subscriber *messenger, call Observer) {
		for {
			if msg, ok := <-subscriber.message; ok {
				// Close if callback returns false.
				select {
				case <-ctx.Done():
					return
				default:
					call(msg)
				}
			}
		}
	}(s, cb)
}
