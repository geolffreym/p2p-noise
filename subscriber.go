package noise

import (
	"context"
)

// Subscriber work as message synchronization.
type Subscriber struct {
	notification chan Message // Message exchange channel
}

func newSubscriber() *Subscriber {
	return &Subscriber{
		notification: make(chan Message),
	}
}

// Emit synchronized message using not-buffered channel.
func (s *Subscriber) Emit(msg Message) {
	s.notification <- msg
}

// Listen and wait for message synchronization from channel.
// When a new message is added to channel buffer the message is proxied to input channel.
func (s *Subscriber) Listen(ctx context.Context, ch chan<- Message) {
	defer close(s.notification)

	for {
		// Close if callback returns false.
		select {
		case <-ctx.Done():
			return
		case msg := <-s.notification:
			ch <- msg // write only channel chan<-
		}
	}
}
