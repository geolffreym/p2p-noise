// Package errors implements custom errors
//
// Refs: https://www.digitalocean.com/community/tutorials/creating-custom-errors-in-go
package errors

import (
	"fmt"
	"net"
)

// Error represents custom errors based on context
type Error struct {
	Context string // Custom error message
	Err     error  // Inherited error from lower level.
}

// Error give string representation of error based on error type
func (e *Error) Error() string {
	switch e.Err.(type) {
	case *net.OpError:
		return fmt.Sprintf("%s -> %v", e.Context, e.Err)
	default:
		return fmt.Sprintf("%s: %v", e.Context, e.Err)
	}
}

// WrapError factory
func WrapErr(err error, context string) error {
	return &Error{
		Err:     err,
		Context: context,
	}
}
