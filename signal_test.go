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

func (m *MockPeer) Socket() Socket {
	return Socket(LOCAL_ADDRESS)
}

func TestType(t *testing.T) {
	event := NewPeerDetected
	payload := []byte(PAYLOAD)
	message := Signal{header{event}, body{payload}, nil}

	if message.Type() != event {
		t.Errorf("expected message with type %v, got %v", event, message.Type())
	}
}

func TestPayload(t *testing.T) {
	event := MessageReceived
	payload := []byte(PAYLOAD)
	message := Signal{header{event}, body{payload}, nil}

	if string(message.Payload()) != string(payload) {
		t.Errorf("expected message with payload %v, got %v", event, message.Type())
	}

}

func TestReply(t *testing.T) {
	event := NewPeerDetected
	payload := []byte(PAYLOAD)
	msg := []byte("hello")

	conn := &mockConn{}
	address := Socket(LOCAL_ADDRESS)
	peer := newPeer(address, conn)

	context := Signal{header{event}, body{payload}, peer}
	sent, _ := context.Reply(msg)

	if sent != len(msg) {
		t.Error("Expected Reply message sent to `Sent` by peer")
	}

}
