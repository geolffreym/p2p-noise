package errors

import (
	"errors"
	"fmt"
	"testing"
)

const MOCK_ADDRESS = "127.0.0.1:2379"

func TestListenError(t *testing.T) {
	customError := "Fail listening"
	err := errors.New(customError)
	output := WrapListen(err, MOCK_ADDRESS)
	expected := fmt.Sprintf("error trying to listen on %s: %v", MOCK_ADDRESS, err)

	if output.Error() != expected {
		t.Errorf("Expected error: %v, got: %v", expected, output)
	}

}

func TestDialError(t *testing.T) {
	customError := "Fail dial"
	err := errors.New(customError)
	output := WrapDial(err, MOCK_ADDRESS)
	expected := fmt.Sprintf("failed dialing to %s: %v", MOCK_ADDRESS, err)

	if output.Error() != expected {
		t.Errorf("Expected error: %v, got: %v", expected, output)
	}

}

func TestBindingError(t *testing.T) {
	customError := "Fail binding"
	err := errors.New(customError)
	output := WrapBinding(err)
	expected := fmt.Sprintf("connection closed or cannot be established: %v", err)

	if output.Error() != expected {
		t.Errorf("Expected error: %v, got: %v", expected, output)
	}
}

func TestCloseError(t *testing.T) {
	customError := "Fail closing connection"
	err := errors.New(customError)
	output := WrapClose(err)
	expected := fmt.Sprintf("error when shutting down connection: %v", err)

	if output.Error() != expected {
		t.Errorf("Expected error: %v, got: %v", expected, output)
	}
}
