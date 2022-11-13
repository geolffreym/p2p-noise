package noise

import (
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
	session := mockSession(&mockConn{}, nil)
	peer := newPeer(session)
	header := header{peer, NewPeerDetected}
	message := Signal{header, PAYLOAD}

	if message.Payload() != PAYLOAD {
		t.Errorf("expected message with payload %v, got %v", event, message.Type())
	}

}
