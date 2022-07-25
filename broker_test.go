package noise

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	event := newBroker()
	subscriber := newSubscriber()
	event.Register(NewPeerDetected, subscriber)
	event.Register(PeerDisconnected, subscriber)
	event.Register(MessageReceived, subscriber)

	registered := []struct {
		name  string
		event Event
	}{{
		name:  "New peer",
		event: NewPeerDetected,
	}, {
		name:  "Peer Disconnected",
		event: PeerDisconnected,
	}, {
		name:  "Message received",
		event: MessageReceived,
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

func TestUnregister(t *testing.T) {
	broker := newBroker()
	subscriber := newSubscriber()
	broker.Register(MessageReceived, subscriber)
	broker.Register(NewPeerDetected, subscriber)
	// Remove self listening from broker events
	success := broker.Unregister(MessageReceived, subscriber)

	if !success {
		t.Errorf("expected success unregister for valid subscriber %v", subscriber)
	}

}

func TestUnregisterExpectedLen(t *testing.T) {
	broker := newBroker()
	subscriber := newSubscriber()
	broker.Register(MessageReceived, subscriber)
	broker.Register(NewPeerDetected, subscriber)
	lenListeningSubscribed := len(broker.topics[MessageReceived])

	// Only NewPeerDetected should be found.
	if lenListeningSubscribed == 2 {
		t.Errorf("expected MessageReceived event unregistered, got %#v events remaining", lenListeningSubscribed)
	}

}

func TestInvalidUnregister(t *testing.T) {
	broker := newBroker()
	subscriber := newSubscriber()
	// Remove self listening from broker events
	success := broker.Unregister(MessageReceived, subscriber)

	if success {
		t.Errorf("expected fail unregister for invalid subscriber %v", subscriber)
	}

}

func TestTopicAdd(t *testing.T) {
	topic := make(topics)
	subscribed := newSubscriber()

	topic.Add(MessageReceived, subscribed)
	topic.Add(PeerDisconnected, subscribed)
	topic.Add(NewPeerDetected, subscribed)

	m, okMsg := topic[MessageReceived]
	p, okPeerDisconnect := topic[PeerDisconnected]
	n, okNewPeer := topic[NewPeerDetected]

	notFoundKeys := !okMsg || !okNewPeer || !okPeerDisconnect
	emptyKeys := len(m) == 0 || len(p) == 0 || len(n) == 0

	if notFoundKeys || emptyKeys {
		t.Errorf("expected topics keys contains added events: MessageReceived, PeerDisconnected, NewPeerDetected")
	}
}

func TestTopicRemove(t *testing.T) {
	topic := make(topics)
	subscribed := newSubscriber()

	topic.Add(MessageReceived, subscribed)
	topic.Add(PeerDisconnected, subscribed)
	removed := topic.Remove(MessageReceived, subscribed)

	emptyKey := len(topic[MessageReceived]) == 0
	integrityCheck := len(topic[PeerDisconnected]) > 0

	// If subscribed not removed and topic with subscribers has entries
	if !removed || !emptyKey || !integrityCheck {
		t.Errorf("expected topics MessageReceived not found after remove")
	}
}

func TestTopicRemoveInvalid(t *testing.T) {
	topic := make(topics)
	subscribed := newSubscriber()

	topic.Add(MessageReceived, subscribed)
	topic.Add(PeerDisconnected, subscribed)
	removed := topic.Remove(NewPeerDetected, subscribed)

	_, existsNewPeer := topic[NewPeerDetected]

	// If subscribed not removed and topic with subscribers has entries
	if removed || existsNewPeer {
		t.Errorf("expected topics NewPeerDetected not found if not registered")
	}
}

func TestPublish(t *testing.T) {
	var result SignalContext
	subscriber := newSubscriber()
	broker := newBroker()

	broker.Register(NewPeerDetected, subscriber)
	message := newSignalContext(NewPeerDetected, []byte("Hello"), nil)

	broker.Publish(message)

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
	broker.Register(NewPeerDetected, subscriber)
	message = newSignalContext(NewPeerDetected, []byte(""), nil)

	broker.Publish(message)
	// Get next message from channel
	// Expected Emit called to set message
	result = <-subscriber.notification
	if result.Type() != NewPeerDetected {
		t.Errorf("expected message type equal to %#v", NewPeerDetected)
	}

}

func TestIndexOf(t *testing.T) {
	slice := []*subscriber{
		newSubscriber(),
		newSubscriber(),
		newSubscriber(),
		newSubscriber(),
		newSubscriber(),
		newSubscriber(),
		newSubscriber(),
	}

	// Table driven test
	// For each expected event
	for _, e := range slice {
		t.Run(fmt.Sprintf("Match for %x", e), func(t *testing.T) {
			match := IndexOf(slice, e)
			if ^match == 0 {
				t.Errorf("expected matched existing elements in slice index: %#v", e)
			}
		})

	}
}

func TestInvalidIndexOf(t *testing.T) {
	slice := []int{1, 2, 3, 4, 6, 7, 8, 9}
	match := IndexOf(slice, 5)

	if ^match != 0 {
		t.Error("Number 5 is not in slice cannot be found")
	}
}
