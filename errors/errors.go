// Package errors implements custom big errors wrapper
//
// The main idea with this package is centralize/proxy errors
// all the modules errors could be declared in this package
// any native error can be proxied using this package.
// !Important using errors in imports could cause a conflict.
// !Important if you need to use errors native package please declare it with alias.
//
// Refs: https://www.digitalocean.com/community/tutorials/creating-custom-errors-in-go
package errors

import (
	"errors"
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

// Alias for error factory
func New(context string) error {
	return errors.New(context)
}

// WrapError factory
func WrapErr(err error, context string) error {
	return &Error{
		Err:     err,
		Context: context,
	}
}
