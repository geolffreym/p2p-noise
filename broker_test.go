package noise

import (
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
			subscribed := s.s[0]           // first element in event subscribed

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
	lenListeningSubscribed := len(broker.topics[MessageReceived].s)

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
	emptyKeys := m.size == 0 || p.size == 0 || n.size == 0

	if notFoundKeys || emptyKeys {
		t.Error("expected topics keys contains added events: MessageReceived, PeerDisconnected, NewPeerDetected")
	}
}

func TestTopicRemove(t *testing.T) {
	topic := make(topics)
	subscribed := newSubscriber()

	topic.Add(MessageReceived, subscribed)
	topic.Add(PeerDisconnected, subscribed)
	removed := topic.Remove(MessageReceived, subscribed)

	emptyKey := len(topic[MessageReceived].s) == 0
	integrityCheck := len(topic[PeerDisconnected].s) > 0

	// If subscribed not removed and topic with subscribers has entries
	if !removed || !emptyKey || !integrityCheck {
		t.Errorf("expected topics MessageReceived not found after remove")
	}
}

func TestTopicAddData(t *testing.T) {
	topic := make(topics)
	subscribed := newSubscriber()

	topic.Add(MessageReceived, subscribed)
	topicLen := topic[MessageReceived].Len()
	subscribers := topic[MessageReceived].Subscribers()

	if topicLen == 0 || subscribers[0] != subscribed {
		t.Error("expected existing topics MessageReceived with data len > 0")
	}
}

func TestTopicDataRemove(t *testing.T) {
	topic := make(topics)
	subscribed := newSubscriber()

	topic.Add(MessageReceived, subscribed)
	topic.Remove(MessageReceived, subscribed)

	topicLen := topic[MessageReceived].Len()
	subscribers := topic[MessageReceived].Subscribers()

	if topicLen > 0 && subscribers[0] != nil {
		t.Errorf("expected topics MessageReceived not found in data after remove")
	}
}

func TestTopicRemoveInvalid(t *testing.T) {
	topic := make(topics)
	subscribed := newSubscriber()

	topic.Add(MessageReceived, subscribed)
	topic.Add(PeerDisconnected, subscribed)
	topic.Add(NewPeerDetected, subscribed)
	// Remove by first time the topic
	topic.Remove(NewPeerDetected, subscribed)
	// trying to remove an already removed subscriber
	removed := topic.Remove(NewPeerDetected, subscribed)

	// If subscribed not removed and topic with subscribers has entries
	if removed {
		t.Error("expected topics NewPeerDetected not found if not registered")
	}
}

func TestPublish(t *testing.T) {
	var result Signal
	subscriber := newSubscriber()
	broker := newBroker()

	header1 := header{newPeer(&mockConn{}), NewPeerDetected}
	signaling := Signal{header1, ""}

	broker.Register(NewPeerDetected, subscriber)
	broker.Publish(signaling)

	// First to finish wins
	// Get first message from channel
	// Expected Emit called to set message
	select {
	case result = <-subscriber.notification:
		if string(result.Payload()) != string(signaling.Payload()) {
			t.Error("expected message equal result")
		}
	case <-time.After(1 * time.Second):
		// Wait 1 second to receive message
		t.Error("expected message received after publish")
		t.FailNow() // If fail receiving messages next test will fail too
	}

	// New message for new topic event
	broker.Register(NewPeerDetected, subscriber)

	header2 := header{newPeer(&mockConn{}), NewPeerDetected}
	signaling = Signal{header2, ""}

	// Number of subscribers notified
	notified := broker.Publish(signaling)
	// Get next message from channel
	// Expected Emit called to set message
	result = <-subscriber.notification
	if result.Type() != NewPeerDetected || notified == 0 {
		t.Errorf("expected message type equal to %#v", NewPeerDetected)
	}

}

func TestInvalidPublish(t *testing.T) {
	broker := newBroker()
	header1 := header{newPeer(&mockConn{}), NewPeerDetected}
	signaling := Signal{header1, ""}

	// Number of subscribers notified
	notified := broker.Publish(signaling)
	// Get next message from channel
	// Expected Emit called to set message
	if notified > 0 {
		t.Error("expected notified to 0 subscribers if not topic registered")
	}

}
