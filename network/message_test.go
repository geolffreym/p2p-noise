package network

import (
	"testing"
)

func TestNewMessage(t *testing.T) {

	event := Event(CLOSED_CONNECTION)
	payload := []byte("hello test")
	message := NewMessage(event, payload, nil)

	if message.Type != event {
		t.Errorf("Expected message with type %v, got %v", event, message.Type)
	}

	if string(message.Payload) != string(payload) {
		t.Errorf("Expected message with payload %v, got %v", event, message.Type)
	}

	if message.Peer != nil {
		t.Errorf("Expected message with nil Peer interface, got %v", message.Peer)
	}

}
