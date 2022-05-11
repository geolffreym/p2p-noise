package errors

import "fmt"

// WrapListen error represent an issue for node address listening
func WrapListen(err error, addr string) error {
	return WrapErr(err, fmt.Sprintf("error trying to listen on %s", addr))
}

// WrapDial error represent an issue trying to dial a node address
func WrapDial(err error, addr string) error {
	return WrapErr(err, fmt.Sprintf("failed dialing to %s", addr))
}

// WrapBinding error represent an issue accepting connections
func WrapBinding(err error) error {
	return WrapErr(err, "connection closed or cannot be established")
}

// WrapClose error represent an issue trying to close connections
func WrapClose(err error) error {
	return WrapErr(err, "error when shutting down connection")
}
