package noise

// subscriber intercept Signal from already subscribed topics in broker
// Handle actions to emit or receive events.
type subscriber struct {
	notification chan Signal // Message exchange channel
}

func newSubscriber() *subscriber {
	return &subscriber{
		make(chan Signal),
	}
}

// Close shuts down the subscriber
func (s *subscriber) Close() {
	close(s.notification)
}

// Emit synchronized message using not-buffered channel.
func (s *subscriber) Emit(msg Signal) {
	s.notification <- msg
}

// Listen and wait for Signal synchronization from channel.
// When a new Signal is added to channel buffer the message is proxied to input channel.
func (s *subscriber) Listen(ch chan<- Signal) {
	for {
		msg := <-s.notification //proxy channel
		ch <- msg               // write only channel chan<-
	}
}
