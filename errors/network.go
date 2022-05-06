package errors

import "fmt"

func Listen(err error, addr string) error {
	return WrapErr(err, fmt.Sprintf("error trying to listen on %s", addr))
}

func Dial(err error, addr string) error {
	return WrapErr(err, fmt.Sprintf("failed to connect to %s", addr))
}

func Binding(err error) error {
	return WrapErr(err, "connection closed or cannot be established")
}

func Close(err error) error {
	return WrapErr(err, "error when shutting down connection")
}
