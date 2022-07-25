package noise

import (
	"context"
)

// Subscriber work as message synchronization.
type subscriber struct {
	notification chan SignalContext // Message exchange channel
}

func newSubscriber() *subscriber {
	return &subscriber{
		make(chan SignalContext),
	}
}

// Emit synchronized message using not-buffered channel.
func (s *subscriber) Emit(msg SignalContext) {
	s.notification <- msg
}

// Listen and wait for message synchronization from channel.
// When a new message is added to channel buffer the message is proxied to input channel.
func (s *subscriber) Listen(ctx context.Context, ch chan<- SignalContext) {
	// Wait until message synchronization finish to close channel
	defer close(s.notification)

	for {
		// Close if callback returns false.
		// select await both of these values simultaneously, executing each one as it arrives.
		select {
		case <-ctx.Done():
			return
		case msg := <-s.notification:
			ch <- msg // write only channel chan<-
		}
	}
}
