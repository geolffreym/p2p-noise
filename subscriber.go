package noise

import (
	"context"
	"sync"
)

// Subscriber synchronize messages from events
// and keep a record of current subscriptions for events.
type Subscriber struct {
	sync.RWMutex              // Mutual exclusion
	notification chan Message // Message exchange channel
}

// Messenger factory
func newSubscriber() *Subscriber {
	return &Subscriber{
		notification: make(chan Message),
	}
}

// Send Message to channel buffer.
func (s *Subscriber) Emit(msg Message) {
	// Lock exclusive writing
	// https://gobyexample.com/mutexes
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()
	s.notification <- msg
}

// Listen and wait fot Message synchronization from channel.
// When a new message is added to channel buffer the
// Observer is executed with new Message propagated as param.
func (s *Subscriber) Listen(ctx context.Context, cb Observer) {
	for {
		if msg, ok := <-s.notification; ok {
			// Close if callback returns false.
			select {
			case <-ctx.Done():
				return
			default:
				cb(msg)
			}
		}
	}
}
