package errors

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
	output := Listening(err, MOCK_ADDRESS)
	expected := fmt.Sprintf("error trying to listen on %s: %v", MOCK_ADDRESS, err)

	if output.Error() != expected {
		t.Errorf(STATEMENT, expected, output)
	}

}

func TestDialingError(t *testing.T) {
	customError := "Fail dial"
	err := errors.New(customError)
	output := Dialing(err, MOCK_ADDRESS)
	expected := fmt.Sprintf("failed dialing to %s: %v", MOCK_ADDRESS, err)

	if output.Error() != expected {
		t.Errorf(STATEMENT, expected, output)
	}

}

func TestBindingError(t *testing.T) {
	customError := "Fail binding"
	err := errors.New(customError)
	output := Binding(err)
	expected := fmt.Sprintf("connection closed or cannot be established: %v", err)

	if output.Error() != expected {
		t.Errorf(STATEMENT, expected, output)
	}
}

func TestClosingError(t *testing.T) {
	customError := "Fail closing connection"
	err := errors.New(customError)
	output := Closing(err)
	expected := fmt.Sprintf("error when shutting down connection: %v", err)

	if output.Error() != expected {
		t.Errorf(STATEMENT, expected, output)
	}
}

func TestExceededError(t *testing.T) {

	output := Exceeded(10)
	expected := fmt.Sprintf("it is not possible to accept more than %d connections: max peers exceeded", 10)

	if output.Error() != expected {
		t.Errorf(STATEMENT, expected, output)
	}
}

func TestMessageError(t *testing.T) {

	output := Message(MOCK_ADDRESS)
	expected := fmt.Sprintf("error trying to send a message to %v: peer disconnected", MOCK_ADDRESS)

	if output.Error() != expected {
		t.Errorf(STATEMENT, expected, output)
	}
}
