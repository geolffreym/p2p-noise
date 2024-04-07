package noise

import "context"

// subscriber intercept Signal from already subscribed topics in broker
// Handle actions to emit or receive events.
type subscriber struct {
	// No, you don't need to close the channel
	// https://stackoverflow.com/questions/8593645/is-it-ok-to-leave-a-channel-open
	notification chan Signal // Message exchange channel
}

func newSubscriber() *subscriber {
	return &subscriber{
		make(chan Signal),
	}
}

// Emit synchronized message using not-buffered channel.
func (s *subscriber) Emit(msg Signal) {
	s.notification <- msg
}

// Listen and wait for Signal synchronization from channel.
// When a new Signal is added to channel buffer the message is proxied to input channel.
func (s *subscriber) Listen(ctx context.Context, ch chan<- Signal) {
	for {
		// Close if callback returns false.
		// select await both of these values simultaneously, executing each one as it arrives.
		select {
		case <-ctx.Done():
			// It's OK to leave a Go channel open forever and never close it.
			// When the channel is no longer used, it will be garbage collected.
			// But "Closing the channel is a control signal on the channel indicating that no more data follows."
			close(ch)
			return
		case msg := <-s.notification:
			ch <- msg // write only channel chan<-
		}
	}
}
