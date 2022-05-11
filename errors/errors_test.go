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

	// Check assertion
	_, ok := wrapper.(*Error)
	if !ok {
		t.Error("Expected 'error' interface implementation")
	}

	if wrapper.Error() != fmt.Sprintf("%s: %v", context, err) {
		t.Error("Expected context and error wrapper to be equal to output")
	}
}
