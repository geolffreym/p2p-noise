package noise

import (
	"context"
)

// SignalContext as message interface to handle network events.
// Each [signal] keep a state holding original header, body and related peer.
type SignalCtx interface {
	Type() Event
	Payload() []byte
	Reply(msg []byte) (int, error)
}

// Subscriber work as message synchronizer.
// Handle actions to emir or receive events.
type subscriber struct {
	notification chan SignalCtx // Message exchange channel
}

func newSubscriber() *subscriber {
	return &subscriber{
		make(chan SignalCtx),
	}
}

// Emit synchronized message using not-buffered channel.
func (s *subscriber) Emit(msg SignalCtx) {
	s.notification <- msg
}

// Listen and wait for message synchronization from channel.
// When a new message is added to channel buffer the message is proxied to input channel.
func (s *subscriber) Listen(ctx context.Context, ch chan<- SignalCtx) {
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
