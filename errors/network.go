package errors

import "fmt"

// Listening error represent an issue for node address listening
func Listening(err error, addr string) error {
	return WrapErr(err, fmt.Sprintf("error trying to listen on %s", addr))
}

// Dialing error represent an issue trying to dial a node address
func Dialing(err error, addr string) error {
	return WrapErr(err, fmt.Sprintf("failed dialing to %s", addr))
}

// Binding error represent an issue accepting connections
func Binding(err error) error {
	return WrapErr(err, "connection closed or cannot be established")
}

// Closing error represent an issue trying to close connections
func Closing(err error) error {
	return WrapErr(err, "error when shutting down connection")
}
