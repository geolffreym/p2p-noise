package noise

import (
	"context"
	"testing"
	"time"
)

func TestSubscriberListen(t *testing.T) {
	sub := newSubscriber()
	body := body{nil}
	header := header{newPeer(PeerA, &mockConn{}), NewPeerDetected}
	signaling := Signal{header, body}

	canceled := make(chan struct{})
	msg := make(chan Signal)
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
		sub.Emit(signaling)
	}()

	// First to finish wins
	select {
	case <-canceled:
		return
	case <-time.After(3 * time.Second):
		// Wait 1 second to receive message
		t.Error("expected canceled listening after emit")
		t.FailNow()
	}

}
