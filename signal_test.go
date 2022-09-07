package noise

import (
	"fmt"
	"testing"
)

const PAYLOAD = "hello test"

func TestType(t *testing.T) {
	event := NewPeerDetected
	message := Signal{header{nil, event}, PAYLOAD}

	if message.Type() != event {
		t.Errorf("expected message with type %v, got %v", event, message.Type())
	}
}

func TestPayload(t *testing.T) {
	event := MessageReceived

	peer := newPeer(&mockConn{})
	header := header{peer, NewPeerDetected}
	message := Signal{header, PAYLOAD}

	fmt.Print(message.Payload())
	if message.Payload() != PAYLOAD {
		t.Errorf("expected message with payload %v, got %v", event, message.Type())
	}

}

func TestReply(t *testing.T) {
	event := NewPeerDetected
	msg := []byte("hello")
	conn := &mockConn{}

	peer := newPeer(conn)
	header := header{peer, event}
	context := Signal{header, PAYLOAD}

	sent, _ := context.Reply(msg)

	if sent != len(msg) {
		t.Error("Expected Reply message sent to `Sent` by peer")
	}

}
