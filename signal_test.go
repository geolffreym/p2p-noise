package noise

import (
	"testing"
)

const PAYLOAD = "hello test"

type MockPeer struct {
	Exchange []byte
}

func (m *MockPeer) Send(data []byte) (n int, err error) {
	m.Exchange = data
	return len(data), nil
}

func (m *MockPeer) Receive(buf []byte) (n int, err error) {
	return len(m.Exchange), nil
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

// TODO add SignalContext test
