package noise

import (
	"reflect"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	event := newBroker()
	subscriber := newSubscriber()
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

	// Table driven test
	// For each expected event
	for _, e := range registered {
		t.Run(e.name, func(t *testing.T) {
			s, ok := event.topics[e.event] // Registered events
			subscribed := s[0]             // first element in event subscribed

			if !ok {
				t.Errorf("expected event %#v, get registered", e)
			}

			if !reflect.DeepEqual(subscribed, subscriber) {
				t.Errorf("expected event subscriber registered equal to %#v, got %#v", subscriber, s)
			}
		})

	}
}

func TestTopicAdd(t *testing.T) {
	topic := make(Topics)
	subscribed := newSubscriber()

	topic.Add(SELF_LISTENING, subscribed)
	topic.Add(CLOSED_CONNECTION, subscribed)
	topic.Add(NEWPEER_DETECTED, subscribed)

	_, okListening := topic[SELF_LISTENING]
	_, okClosed := topic[CLOSED_CONNECTION]
	_, okNewPeer := topic[NEWPEER_DETECTED]

	if !okListening || !okNewPeer || !okClosed {
		t.Errorf("expected topics keys contains added events")
	}
}

func TestPublish(t *testing.T) {
	var result Message
	subscriber := newSubscriber()
	event := newBroker()

	event.Register(SELF_LISTENING, subscriber)
	message := newMessage(SELF_LISTENING, []byte("hello test 1"))
	event.Publish(message)

	// First to finish wins
	// Get first message from channel
	// Expected Emit called to set message
	select {
	case result = <-subscriber.notification:
		if string(result.Payload()) != string(message.Payload()) {
			t.Errorf("expected message equal result")
		}
	case <-time.After(1 * time.Second):
		// Wait 1 second to receive message
		t.Errorf("expected message received after publish")
		t.FailNow() // If fail receiving messages next test will fail too
	}

	// New message for new topic event
	event.Register(NEWPEER_DETECTED, subscriber)
	message = newMessage(NEWPEER_DETECTED, []byte(""))
	event.Publish(message)

	// Get next message from channel
	// Expected Emit called to set message
	result = <-subscriber.notification
	if result.Type() != NEWPEER_DETECTED {
		t.Errorf("expected message type equal to %#v", NEWPEER_DETECTED)
	}

}
