package network

import (
	"reflect"
	"testing"
)

func TestRegister(t *testing.T) {
	subscriber := NewSubscriber()
	event := make(Events)
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
			s, ok := event[e.event] // Registered events
			subscribed := s[0]      // first element in event subscribed

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
	event := make(Events)
	event.Register(SELF_LISTENING, subscriber)
	message := NewMessage(SELF_LISTENING, []byte("hello test 1"), nil)
	event.Publish(message)

	// Get first message from channel
	// Expected Emit called to set message
	result := <-subscriber.message
	if string(result.Payload) != string(message.Payload) {
		t.Errorf("Expected message equal result")
	}

	// Modified message to publish and runtime listening
	event.Register(NEWPEER_DETECTED, subscriber)
	message.Type = NEWPEER_DETECTED
	message.Payload = []byte("hello test 2")
	event.Publish(message)

	// Get next message from channel
	// Expected Emit called to set message
	result = <-subscriber.message
	if string(result.Payload) != string(message.Payload) {
		t.Errorf("Expected message equal result")
	}

}
