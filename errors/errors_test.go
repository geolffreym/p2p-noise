package errors

import (
	"errors"
	"fmt"
	"testing"
)

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
