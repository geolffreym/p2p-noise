package noise

import (
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
// ref: https://go.dev/ref/spec#Conversions
// https://pkg.go.dev/unsafe#Pointer
func newIDFromString(s string) ID {
	// "no-copy" convert to ID from string.
	// If the type starts with the operator * or <-, it must be parenthesized when necessary to avoid ambiguity.
	// 1- unsafe.Pointer(&s) <- create a pointer from string address
	// 2- (*ID)(unsafe.Pointer(&s)) <- cast the pointer to *ID pointer
	// 3 - *(*ID)(unsafe.Pointer(&s)) <- get the value in memory address
	// TODO could by replaced by unsafe.StringData(s) in go >=1.20
	return *(*ID)(unsafe.Pointer(&s))
}

// newBlake2ID creates a new id blake2 hash based.
func newBlake2ID(plaintext []byte) ID {
	var id ID
	hash := blake2(plaintext)
	copy(id[:], hash)
	return id
}
