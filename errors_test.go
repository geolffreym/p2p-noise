package noise

import (
	"errors"
	"fmt"
	"testing"
)

const MOCK_ADDRESS = "127.0.0.1:2379"
const STATEMENT = "Expected error: %v, got: %v"

func TestListeningError(t *testing.T) {
	customError := "Fail listening"
	err := errors.New(customError)
	output := ErrSelfListening(err, MOCK_ADDRESS)
	expected := fmt.Sprintf("net: error trying to listen on %s -> %v", MOCK_ADDRESS, err)

	if output.Error() != expected {
		t.Errorf(STATEMENT, expected, output)
	}

}

func TestDialingError(t *testing.T) {
	customError := "Fail dial"
	err := errors.New(customError)
	output := ErrDialingNode(err, MOCK_ADDRESS)
	expected := fmt.Sprintf("net: failed dialing to %s -> %v", MOCK_ADDRESS, err)

	if output.Error() != expected {
		t.Errorf(STATEMENT, expected, output)
	}

}

func TestBindingError(t *testing.T) {
	customError := "Fail binding"
	err := errors.New(customError)
	output := ErrBindingConnection(err)
	expected := fmt.Sprintf("net: connection closed or cannot be established -> %v", err)

	if output.Error() != expected {
		t.Errorf(STATEMENT, expected, output)
	}
}

func TestClosingError(t *testing.T) {
	customError := "Fail closing connection"
	err := errors.New(customError)
	output := ErrClosingConnection(err)
	expected := fmt.Sprintf("net: error when shutting down connection -> %v", err)

	if output.Error() != expected {
		t.Errorf(STATEMENT, expected, output)
	}
}

func TestExceededMaxPeersError(t *testing.T) {

	output := ErrExceededMaxPeers(10)
	expected := fmt.Sprintf("overflow: it is not possible to accept more than %d connections -> max peers exceeded", 10)

	if output.Error() != expected {
		t.Errorf(STATEMENT, expected, output)
	}
}

func TestExceededMaxPayloadSize(t *testing.T) {

	output := ErrExceededMaxPayloadSize(10 << 20)
	expected := fmt.Sprintf("overflow: it is not possible to accept more than %d bytes -> max payload size exceeded", 10<<20)

	if output.Error() != expected {
		t.Errorf(STATEMENT, expected, output)
	}
}

func TestSendingMessageToInvalidPeer(t *testing.T) {

	output := ErrSendingMessageToInvalidPeer(MOCK_ADDRESS)
	expected := fmt.Sprintf("ops: error trying to send a message to %v -> peer disconnected", MOCK_ADDRESS)

	if output.Error() != expected {
		t.Errorf(STATEMENT, expected, output)
	}
}
