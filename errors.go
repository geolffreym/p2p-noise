package noise

import (
	"errors"
	"fmt"
)

// [NetError] represents errors related to network communication.
type NetError struct {
	Context string
	Err     error
}

// Error give string representation of error based on error type.
func (e NetError) Error() string {
	return fmt.Sprintf("net: %s -> %v", e.Context, e.Err)
}

// [OperationalError] represents an error that occurred when an operation in node failed.
// eg. Send a new message to invalid or not connected peer.
// eg. Error during Handshake.
type OperationalError struct {
	Context string
	Err     error
}

// Error give string representation of error based on error type.
func (e OperationalError) Error() string {
	return fmt.Sprintf("ops: %s -> %v", e.Context, e.Err)
}

// [OverflowError] error represents a problem with the maximum setting of a parameter being exceeded.
// eg. MaxPeersConnected exceeded for incoming connections.
type OverflowError struct {
	Context string
	Err     error
}

// Error give string representation of error based on error type.
func (e OverflowError) Error() string {
	return fmt.Sprintf("overflow: %s -> %v", e.Context, e.Err)
}

// errDialingNode error represent an issue trying to dial a node address.
func errDialingNode(err error) error {
	return &NetError{"error during dialing", err}
}

// errBindingConnection error represent an issue accepting connections.
func errBindingConnection(err error) error {
	return &NetError{"connection closed or cannot be established", err}
}

// errExceededMaxPeers error represent an issue if number of active connections exceed max peer connected.
func errExceededMaxPeers(max uint8) error {
	return &OverflowError{
		fmt.Sprintf("it is not possible to accept more than %d connections", max),
		errors.New("max peers exceeded"),
	}
}

// errExceededMaxPayloadSize error represent an issue if received message size exceed max payload size.
func errExceededMaxPayloadSize(max uint32) error {
	return &OverflowError{
		fmt.Sprintf("it is not possible to accept more than %d bytes", max),
		errors.New("max payload size exceeded"),
	}
}

// errSendingMessage error represent an issue trying to send a message.
func errSendingMessage(err error) error {
	return &OperationalError{"error sending message", err}
}

// errDuringHandshake error represent an issue during handshake with peer.
func errDuringHandshake(err error) error {
	return &OperationalError{"error during handshake", err}
}
