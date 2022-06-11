package noise

import (
	"reflect"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	event := newBroker()
	subscriber := newSubscriber()
	event.Register(SelfListening, subscriber)
	event.Register(NewPeerDetected, subscriber)
	event.Register(ClosedConnection, subscriber)
	event.Register(MessageReceived, subscriber)

	registered := []struct {
		name  string
		event Event
	}{{
		name:  "Listening",
		event: SelfListening,
	}, {
		name:  "New peer",
		event: NewPeerDetected,
	}, {
		name:  "Closed connection",
		event: ClosedConnection,
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
	broker.Register(SelfListening, subscriber)
	broker.Register(NewPeerDetected, subscriber)
	// Remove self listening from broker events
	success := broker.Unregister(SelfListening, subscriber)
	lenListeningSubscribed := len(broker.topics[SelfListening])

	if !success {
		t.Errorf("expected success unregister for valid subscriber %v", subscriber)
	}
	// Only NewPeerDetected should be found.
	if lenListeningSubscribed > 0 {
		t.Errorf("expected SelfListening event unregistered, got %#v events remaining", lenListeningSubscribed)
	}

}

func TestInvalidUnregister(t *testing.T) {
	broker := newBroker()
	subscriber := newSubscriber()
	// Remove self listening from broker events
	success := broker.Unregister(SelfListening, subscriber)

	if success {
		t.Errorf("expected fail unregister for invalid subscriber %v", subscriber)
	}

}

func TestTopicAdd(t *testing.T) {
	topic := make(Topics)
	subscribed := newSubscriber()

	topic.Add(SelfListening, subscribed)
	topic.Add(ClosedConnection, subscribed)
	topic.Add(NewPeerDetected, subscribed)

	_, okListening := topic[SelfListening]
	_, okClosed := topic[ClosedConnection]
	_, okNewPeer := topic[NewPeerDetected]

	if !okListening || !okNewPeer || !okClosed {
		t.Errorf("expected topics keys contains added events")
	}
}

func TestPublish(t *testing.T) {
	var result Message
	subscriber := newSubscriber()
	broker := newBroker()

	broker.Register(SelfListening, subscriber)
	message := newMessage(SelfListening, []byte("hello test 1"))
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
	message = newMessage(NewPeerDetected, []byte(""))
	broker.Publish(message)

	// Get next message from channel
	// Expected Emit called to set message
	result = <-subscriber.notification
	if result.Type() != NewPeerDetected {
		t.Errorf("expected message type equal to %#v", NewPeerDetected)
	}

}
