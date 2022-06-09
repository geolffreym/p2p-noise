package errors

import (
	"errors"
	"fmt"
	"testing"

	// Avoid renaming imports except to avoid a name collision; good package names should not require renaming.
	// In the event of collision, prefer to rename the most local or project-specific import.
	errors_ "github.com/geolffreym/p2p-noise/errors"
)

func TestWrapError(t *testing.T) {
	err := errors.New("wrap test")
	context := "testing errors"
	wrapper := errors_.WrapErr(err, context)
	expected := fmt.Sprintf("%s: %v", context, err)

	// Check assertion
	_, ok := wrapper.(*errors_.Error)
	if !ok {
		t.Error("Expected 'error' interface implementation")
	}

	if wrapper.Error() != expected {
		t.Error("Expected context and error wrapper to be equal to output")
	}
}
