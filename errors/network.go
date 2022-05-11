package errors

import "fmt"

// Listen error represent an issue for node address listening
func Listen(err error, addr string) error {
	return WrapErr(err, fmt.Sprintf("error trying to listen on %s", addr))
}

// Dial error represent an issue trying to dial a node address
func Dial(err error, addr string) error {
	return WrapErr(err, fmt.Sprintf("failed dialing to %s", addr))
}

// Binding error represent an issue accepting connections
func Binding(err error) error {
	return WrapErr(err, "connection closed or cannot be established")
}

// Close error represent an issue trying to close connections
func Close(err error) error {
	return WrapErr(err, "error when shutting down connection")
}
