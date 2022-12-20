package noise

import (
	"reflect"
	"unsafe"
)

// [ID] it's identity provider for peer.
type ID [32]byte

// Bytes return a byte slice representation for id.
func (i ID) Bytes() []byte {
	return i[:]
}

// String return a string representation for 32-bytes hash.
func (i ID) String() string {
	return (string)(i[:])
}

// newIDFromString creates a new ID from string.
// ref: https://stackoverflow.com/questions/59209493/how-to-use-unsafe-get-a-byte-slice-from-a-string-without-memory-copy
func newIDFromString(s string) ID {
	// "no-copy" convert to ID from string.
	return *(*ID)(unsafe.Pointer(
		(*reflect.StringHeader)(unsafe.Pointer(&s)).Data),
	)
}

// newBlake2ID creates a new id blake2 hash based.
func newBlake2ID(plaintext []byte) ID {
	var id ID
	// Hash
	hash := blake2(plaintext)
	// Populate id
	copy(id[:], hash)
	return id
}
