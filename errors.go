package noise

import (
	"errors"
	"fmt"
)

// NetError represents custom errors based on context
type NetError struct {
	Context string // Custom error message
	Err     error  // Inherited error from lower level.
}

// Error give string representation of error based on error type.
func (e NetError) Error() string {
	return fmt.Sprintf("net: %s -> %v", e.Context, e.Err)
}

type OperationalError struct {
	Context string // Custom error message
	Err     error  // Inherited error from lower level.
}

// Error give string representation of error based on error type.
func (e OperationalError) Error() string {
	return fmt.Sprintf("ops: %s -> %v", e.Context, e.Err)
}

type OverflowError struct {
	Context string // Custom error message
	Err     error  // Inherited error from lower level.
}

// Error give string representation of error based on error type.
func (e OverflowError) Error() string {
	return fmt.Sprintf("overflow: %s -> %v", e.Context, e.Err)
}

// errSelfListening error represent an issue for node address listening.
func errSelfListening(err error, addr string) error {
	return &NetError{fmt.Sprintf("error trying to listen on %s", addr), err}
}

// errDialingNode error represent an issue trying to dial a node address.
func errDialingNode(err error, addr string) error {
	return &NetError{fmt.Sprintf("failed dialing to %s", addr), err}
}

// errBindingConnection error represent an issue accepting connections.
func errBindingConnection(err error) error {
	return &NetError{"connection closed or cannot be established", err}
}

// errClosingConnection error represent an issue trying to close connections.
func errClosingConnection(err error) error {
	return &NetError{"error when shutting down connection", err}
}

// errExceededMaxPeers error represent an issue if number of active connections exceed max peer connected.
func errExceededMaxPeers(max uint8) error {
	return &OverflowError{
		fmt.Sprintf("it is not possible to accept more than %d connections", max),
		errors.New("max peers exceeded"),
	}
}

// ErrExceededMaxPeers error represent an issue if number of active connections exceed max peer connected.
func errExceededMaxPayloadSize(max uint32) error {
	return &OverflowError{
		fmt.Sprintf("it is not possible to accept more than %d bytes", max),
		errors.New("max payload size exceeded"),
	}
}

func errSendingMessageToInvalidPeer(addr string) error {
	return &OperationalError{
		fmt.Sprintf("error trying to send a message to %s", addr),
		errors.New("peer disconnected"),
	}
}
