package errors

import (
	"fmt"
	"testing"

	errors "github.com/geolffreym/p2p-noise/errors"
)

func TestWrapError(t *testing.T) {
	err := errors.New("wrap test")
	context := "testing errors"
	wrapper := errors.WrapErr(err, context)
	expected := fmt.Sprintf("%s: %v", context, err)

	// Check assertion
	_, ok := wrapper.(*errors.Error)
	if !ok {
		t.Error("Expected 'error' interface implementation")
	}

	if wrapper.Error() != expected {
		t.Error("Expected context and error wrapper to be equal to output")
	}
}
