package noise

import (
	"context"
	"testing"
	"time"
)

func TestListen(t *testing.T) {
	sub := newSubscriber()
	message := newSignalContext(SelfListening, []byte("hello test 1"), nil)

	canceled := make(chan struct{})
	msg := make(chan SignalContext)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		// wait for new message emitted then cancel listening
		<-msg
		cancel()
	}()

	go func() {
		// after stop listening loop expect trigger canceled
		sub.Listen(ctx, msg)
		canceled <- struct{}{}
	}()

	go func() {
		// send message after 1 second
		time.Sleep(1 * time.Second)
		sub.Emit(message)
	}()

	// First to finish wins
	select {
	case <-canceled:
		t.Log("canceled channel")
	case <-time.After(3 * time.Second):
		// Wait 1 second to receive message
		t.Errorf("expected canceled listening after emit")
		t.FailNow()
	}

}
