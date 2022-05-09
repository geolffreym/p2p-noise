// Package errors implements custom errors
//
// Refs: Based on https://www.digitalocean.com/community/tutorials/creating-custom-errors-in-go
package errors

import (
	"fmt"
)

type Error struct {
	Context string
	Err     error
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %v", e.Context, e.Err)
}

// Error factory
func WrapErr(err error, context string) error {
	return &Error{
		Err:     err,
		Context: context,
	}
}
