package noise

import (
	"reflect"
	"unsafe"
)

// [ID] serves as the identity for peers.
// It facilitates addressability in router table.
type ID [32]byte

// Bytes return a byte slice representation for id.
func (i ID) Bytes() []byte {
	return i[:]
}

// String return a string representation for 32-bytes hash.
// ref: https://go.dev/ref/spec#Conversions
func (i ID) String() string {
	return string(i[:])
}

// newIDFromString creates a new ID from string.
// ref: https://stackoverflow.com/questions/59209493/how-to-use-unsafe-get-a-byte-slice-from-a-string-without-memory-copy
// ref: https://go.dev/ref/spec#Conversions
func newIDFromString(s string) ID {
	// "no-copy" convert to ID from string.
	// If the type starts with the operator * or <-, it must be parenthesized when necessary to avoid ambiguity.
	return *(*ID)(unsafe.Pointer(
		(*reflect.StringHeader)(unsafe.Pointer(&s)).Data),
	)
}

// newBlake2ID creates a new id blake2 hash based.
func newBlake2ID(plaintext []byte) ID {
	var id ID
	hash := blake2(plaintext)
	copy(id[:], hash)
	return id
}
