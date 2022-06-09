package errors_test

import (
	"errors"
	"fmt"
	"testing"

	// Avoid renaming imports except to avoid a name collision; good package names should not require renaming.
	// In the event of collision, prefer to rename the most local or project-specific import.
	//
	// ref: https://github.com/golang/go/wiki/CodeReviewComments
	errors_ "github.com/geolffreym/p2p-noise/errors"
)

const MOCK_ADDRESS = "127.0.0.1:2379"

func TestListeningError(t *testing.T) {
	customError := "Fail listening"
	err := errors.New(customError)
	output := errors_.Listening(err, MOCK_ADDRESS)
	expected := fmt.Sprintf("error trying to listen on %s: %v", MOCK_ADDRESS, err)

	if output.Error() != expected {
		t.Errorf("Expected error: %v, got: %v", expected, output)
	}

}

func TestDialingError(t *testing.T) {
	customError := "Fail dial"
	err := errors.New(customError)
	output := errors_.Dialing(err, MOCK_ADDRESS)
	expected := fmt.Sprintf("failed dialing to %s: %v", MOCK_ADDRESS, err)

	if output.Error() != expected {
		t.Errorf("Expected error: %v, got: %v", expected, output)
	}

}

func TestBindingError(t *testing.T) {
	customError := "Fail binding"
	err := errors.New(customError)
	output := errors_.Binding(err)
	expected := fmt.Sprintf("connection closed or cannot be established: %v", err)

	if output.Error() != expected {
		t.Errorf("Expected error: %v, got: %v", expected, output)
	}
}

func TestClosingError(t *testing.T) {
	customError := "Fail closing connection"
	err := errors.New(customError)
	output := errors_.Closing(err)
	expected := fmt.Sprintf("error when shutting down connection: %v", err)

	if output.Error() != expected {
		t.Errorf("Expected error: %v, got: %v", expected, output)
	}
}
