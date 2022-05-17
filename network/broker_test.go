package network

import (
	"reflect"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	subscriber := NewSubscriber()
	event := NewEvents()
	event.Register(SELF_LISTENING, subscriber)
	event.Register(NEWPEER_DETECTED, subscriber)
	event.Register(CLOSED_CONNECTION, subscriber)
	event.Register(MESSAGE_RECEIVED, subscriber)

	registered := []struct {
		name  string
		event Event
	}{{
		name:  "Listening",
		event: SELF_LISTENING,
	}, {
		name:  "New peer",
		event: NEWPEER_DETECTED,
	}, {
		name:  "Closed connection",
		event: CLOSED_CONNECTION,
	}, {
		name:  "Message received",
		event: MESSAGE_RECEIVED,
	}}

	// For each expected event
	for _, e := range registered {
		t.Run(e.name, func(t *testing.T) {
			s, ok := event.topics[e.event] // Registered events
			subscribed := s[0]             // first element in event subscribed

			if !ok {
				t.Errorf("Expected event %#v, get registered", e)
			}

			if !reflect.DeepEqual(subscribed, subscriber) {
				t.Errorf("Expected event subscriber registered equal to %#v, got %#v", subscriber, s)
			}
		})

	}
}

func TestPublish(t *testing.T) {
	subscriber := NewSubscriber()
	event := NewEvents()
	event.Register(SELF_LISTENING, subscriber)
	message := NewMessage(SELF_LISTENING, []byte("hello test 1"), nil)
	event.Publish(message)
	var result *Message

	// Get first message from channel
	// Expected Emit called to set message
	select {
	case result := <-subscriber.message:
		if string(result.Payload) != string(message.Payload) {
			t.Errorf("Expected message equal result")
		}
	case <-time.After(1 * time.Second):
		// Wait 1 second to receive message
		t.Errorf("Expected message received after publish")
	}

	// New message for new topic event
	event.Register(NEWPEER_DETECTED, subscriber)
	message = NewMessage(NEWPEER_DETECTED, []byte(""), nil)
	event.Publish(message)

	// Get next message from channel
	// Expected Emit called to set message
	result = <-subscriber.message
	if result.Type != Event(NEWPEER_DETECTED) {
		t.Errorf("Expected message type equal to %#v", NEWPEER_DETECTED)
	}

}
