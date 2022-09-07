package noise

import (
	"fmt"
	"testing"
)

const PAYLOAD = "hello test"

type MockPeer struct {
}

func (m *MockPeer) Send(msg []byte) (int, error) {
	return len(msg), nil
}

func (m *MockPeer) Socket() Socket {
	return Socket(LOCAL_ADDRESS)
}

func TestType(t *testing.T) {
	event := NewPeerDetected
	payload := []byte(PAYLOAD)
	message := Signal{header{nil, event}, payload}

	if message.Type() != event {
		t.Errorf("expected message with type %v, got %v", event, message.Type())
	}
}

func TestPayload(t *testing.T) {
	event := MessageReceived

	body := []byte(PAYLOAD)
	peer := newPeer(PeerA, &mockConn{})
	header := header{peer, NewPeerDetected}
	message := Signal{header, body}

	fmt.Print(message.Payload())
	if message.Payload() != PAYLOAD {
		t.Errorf("expected message with payload %v, got %v", event, message.Type())
	}

}

func TestReply(t *testing.T) {
	event := NewPeerDetected
	msg := []byte("hello")
	conn := &mockConn{}

	address := Socket(LOCAL_ADDRESS)
	peer := newPeer(address, conn)

	body := []byte(PAYLOAD)
	header := header{peer, event}
	context := Signal{header, body}

	sent, _ := context.Reply(msg)

	if sent != len(msg) {
		t.Error("Expected Reply message sent to `Sent` by peer")
	}

}
