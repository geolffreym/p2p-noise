package network_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/geolffreym/p2p-noise/network"
)

func TestRegister(t *testing.T) {
	event := network.NewEvents()
	subscriber := network.NewMessenger()
	event.Register(network.SELF_LISTENING, subscriber)
	event.Register(network.NEWPEER_DETECTED, subscriber)
	event.Register(network.CLOSED_CONNECTION, subscriber)
	event.Register(network.MESSAGE_RECEIVED, subscriber)

	registered := []struct {
		name  string
		event network.Event
	}{{
		name:  "Listening",
		event: network.SELF_LISTENING,
	}, {
		name:  "New peer",
		event: network.NEWPEER_DETECTED,
	}, {
		name:  "Closed connection",
		event: network.CLOSED_CONNECTION,
	}, {
		name:  "Message received",
		event: network.MESSAGE_RECEIVED,
	}}

	// Table driven test
	// For each expected event
	for _, e := range registered {
		t.Run(e.name, func(t *testing.T) {
			s, ok := event.Topics()[e.event] // Registered events
			subscribed := s[0]               // first element in event subscribed

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
	topic := make(network.Topics)
	subscribed := network.NewMessenger()

	topic.Add(network.SELF_LISTENING, subscribed)
	topic.Add(network.CLOSED_CONNECTION, subscribed)
	topic.Add(network.NEWPEER_DETECTED, subscribed)

	_, okListening := topic[network.SELF_LISTENING]
	_, okClosed := topic[network.CLOSED_CONNECTION]
	_, okNewPeer := topic[network.NEWPEER_DETECTED]

	if !okListening || !okNewPeer || !okClosed {
		t.Errorf("expected topics keys contains added events")
	}
}

func TestPublish(t *testing.T) {
	var result network.Message
	subscriber := network.NewMessenger()
	event := network.NewEvents()

	event.Register(network.SELF_LISTENING, subscriber)
	message := network.NewMessage(network.SELF_LISTENING, []byte("hello test 1"), nil)
	event.Publish(message)

	// First to finish wins
	// Get first message from channel
	// Expected Emit called to set message
	select {
	case result = <-subscriber.Message():
		if string(result.Payload()) != string(message.Payload()) {
			t.Errorf("expected message equal result")
		}
	case <-time.After(1 * time.Second):
		// Wait 1 second to receive message
		t.Errorf("expected message received after publish")
		t.FailNow() // If fail receiving messages next test will fail too
	}

	// New message for new topic event
	event.Register(network.NEWPEER_DETECTED, subscriber)
	message = network.NewMessage(network.NEWPEER_DETECTED, []byte(""), nil)
	event.Publish(message)

	// Get next message from channel
	// Expected Emit called to set message
	result = <-subscriber.Message()
	if result.Type() != network.Event(network.NEWPEER_DETECTED) {
		t.Errorf("expected message type equal to %#v", network.NEWPEER_DETECTED)
	}

}
