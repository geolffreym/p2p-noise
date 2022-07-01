package noise

import (
	"errors"
	"fmt"
	"testing"
)

const MOCK_ADDRESS = "127.0.0.1:2379"
const STATEMENT = "Expected error: %v, got: %v"

func TestWrapError(t *testing.T) {
	err := errors.New("wrap test")
	context := "testing errors"
	wrapper := WrapErr(err, context)
	expected := fmt.Sprintf("%s: %v", context, err)

	// Check assertion
	_, ok := wrapper.(*Error)
	if !ok {
		t.Error("Expected 'error' interface implementation")
	}

	if wrapper.Error() != expected {
		t.Error("Expected context and error wrapper to be equal to output")
	}
}

func TestListeningError(t *testing.T) {
	customError := "Fail listening"
	err := errors.New(customError)
	output := ErrSelfListening(err, MOCK_ADDRESS)
	expected := fmt.Sprintf("error trying to listen on %s: %v", MOCK_ADDRESS, err)

	if output.Error() != expected {
		t.Errorf(STATEMENT, expected, output)
	}

}

func TestDialingError(t *testing.T) {
	customError := "Fail dial"
	err := errors.New(customError)
	output := ErrDialingNode(err, MOCK_ADDRESS)
	expected := fmt.Sprintf("failed dialing to %s: %v", MOCK_ADDRESS, err)

	if output.Error() != expected {
		t.Errorf(STATEMENT, expected, output)
	}

}

func TestBindingError(t *testing.T) {
	customError := "Fail binding"
	err := errors.New(customError)
	output := ErrBindingConnection(err)
	expected := fmt.Sprintf("connection closed or cannot be established: %v", err)

	if output.Error() != expected {
		t.Errorf(STATEMENT, expected, output)
	}
}

func TestClosingError(t *testing.T) {
	customError := "Fail closing connection"
	err := errors.New(customError)
	output := ErrClosingConnection(err)
	expected := fmt.Sprintf("error when shutting down connection: %v", err)

	if output.Error() != expected {
		t.Errorf(STATEMENT, expected, output)
	}
}

func TestExceededError(t *testing.T) {

	output := ErrExceededMaxPeers(10)
	expected := fmt.Sprintf("it is not possible to accept more than %d connections: max peers exceeded", 10)

	if output.Error() != expected {
		t.Errorf(STATEMENT, expected, output)
	}
}

func TestMessageError(t *testing.T) {

	output := ErrSendingMessage(MOCK_ADDRESS)
	expected := fmt.Sprintf("error trying to send a message to %v: peer disconnected", MOCK_ADDRESS)

	if output.Error() != expected {
		t.Errorf(STATEMENT, expected, output)
	}
}
