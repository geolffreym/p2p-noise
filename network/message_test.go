package network

import (
	"reflect"
	"testing"
)

func TestNewMessage(t *testing.T) {

	event := Event(CLOSED_CONNECTION)
	payload := []byte("hello test")
	message := NewMessage(event, payload, nil)

	if reflect.TypeOf(message) != reflect.TypeOf((Message)(nil)) {
		t.Errorf("expected *Message, got %#v", message)
		t.FailNow() // If fail type PeerImp assertion next tests will fail too
	}

	if message.Type() != event {
		t.Errorf("Expected message with type %v, got %v", event, message.Type())
	}

	if string(message.Payload()) != string(payload) {
		t.Errorf("Expected message with payload %v, got %v", event, message.Type())
	}

	if message.Peer() != nil {
		t.Errorf("Expected message with nil Peer interface, got %v", message.Peer())
	}

}
