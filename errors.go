package noise

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

// Error give string representation of error based on error type.
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

// ErrSelfListening error represent an issue for node address listening.
func ErrSelfListening(err error, addr string) error {
	return WrapErr(err, fmt.Sprintf("error trying to listen on %s", addr))
}

// ErrDialingNode error represent an issue trying to dial a node address.
func ErrDialingNode(err error, addr string) error {
	return WrapErr(err, fmt.Sprintf("failed dialing to %s", addr))
}

// ErrBindingConnection error represent an issue accepting connections.
func ErrBindingConnection(err error) error {
	return WrapErr(err, "connection closed or cannot be established")
}

// ErrClosingConnection error represent an issue trying to close connections.
func ErrClosingConnection(err error) error {
	return WrapErr(err, "error when shutting down connection")
}

// ErrExceededMaxPeers error represent an issue if number of active connections exceed max peer connected.
func ErrExceededMaxPeers(max uint8) error {
	return WrapErr(
		errors.New("max peers exceeded"),
		fmt.Sprintf("it is not possible to accept more than %d connections", max),
	)
}

func ErrSendingMessage(addr string) error {
	return WrapErr(
		errors.New("peer disconnected"),
		fmt.Sprintf("error trying to send a message to %v", addr),
	)
}
