package network_test

import (
	"reflect"
	"testing"

	"github.com/geolffreym/p2p-noise/network"
)

type MockPeer struct {
	Exchange []byte
}

func (m *MockPeer) Send(data []byte) (n int, err error) {
	m.Exchange = data
	return len(data), nil
}

func (m *MockPeer) Receive(buf []byte) (n int, err error) {
	return len(m.Exchange), nil
}

func TestType(t *testing.T) {
	peer := &MockPeer{}
	event := network.Event(network.CLOSED_CONNECTION)
	payload := []byte("hello test")
	message := network.NewMessage(event, payload, peer)

	if message.Type() != event {
		t.Errorf("expected message with type %v, got %v", event, message.Type())
	}
}

func TestPayload(t *testing.T) {
	peer := &MockPeer{}
	event := network.Event(network.CLOSED_CONNECTION)
	payload := []byte("hello test")
	message := network.NewMessage(event, payload, peer)

	if string(message.Payload()) != string(payload) {
		t.Errorf("expected message with payload %v, got %v", event, message.Type())
	}

}

func TestPeer(t *testing.T) {
	peer := &MockPeer{}
	event := network.Event(network.CLOSED_CONNECTION)
	payload := []byte("hello test")
	message := network.NewMessage(event, payload, peer)

	if !reflect.DeepEqual(message.Peer(), peer) {
		t.Errorf("expected message with nil Peer interface, got %v", message.Peer())
	}

}

func TestReply(t *testing.T) {
	peer := &MockPeer{}
	event := network.Event(network.CLOSED_CONNECTION)
	payload := []byte("hello test")
	message := network.NewMessage(event, payload, peer)

	replied := []byte("hey")
	message.Reply(replied)
	mock, _ := message.Peer().(*MockPeer)

	if r := mock.Exchange; len(r) != len(replied) {
		t.Errorf("expected message replied, got %v", mock.Exchange)
	}

}
