package noise

import (
	"testing"
)

const PAYLOAD = "hello test"

type MockPeer struct {
}

func (m *MockPeer) Send(msg []byte) (int, error) {
	return len(msg), nil
}

func TestType(t *testing.T) {
	event := NewPeerDetected
	payload := []byte(PAYLOAD)
	message := signal{
		event,
		payload,
	}

	if message.Type() != event {
		t.Errorf("expected message with type %v, got %v", event, message.Type())
	}
}

func TestPayload(t *testing.T) {
	event := MessageReceived
	payload := []byte(PAYLOAD)
	message := signal{
		event,
		payload,
	}

	if string(message.Payload()) != string(payload) {
		t.Errorf("expected message with payload %v, got %v", event, message.Type())
	}

}

func TestReply(t *testing.T) {
	payload := []byte(PAYLOAD)
	msg := []byte("hello")
	peer := &MockPeer{}

	context := newSignalContext(NewPeerDetected, payload, peer)
	sent, _ := context.Reply(msg)

	if sent != len(msg) {
		t.Error("Expected Reply same as Sent by peer")
	}

}
