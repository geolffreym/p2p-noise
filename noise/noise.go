package noise

import (
	"hash"

	"github.com/flynn/noise"
	"golang.org/x/crypto/blake2b"
)

// BLAKE2 is a cryptographic hash function faster than MD5, SHA-1, SHA-2, and SHA-3.
// ref: https://www.blake2.net/
type blake2Fn struct{}

func (h blake2Fn) HashName() string { return "BLAKE2b" }
func (h blake2Fn) Hash() hash.Hash {
	// New returns a new hash.Hash computing the BLAKE2b checksum with a custom length.
	// A non-nil key turns the hash into a MAC. The key must be between zero and 64 bytes long.
	// The hash size can be a value between 1 and 64 but it is highly recommended to use
	// values equal or greater than:
	// - 32 if BLAKE2b is used as a hash function (The key is zero bytes long).
	// - 16 if BLAKE2b is used as a MAC function (The key is at least 16 bytes long).
	// When the key is nil, the returned hash.Hash implements BinaryMarshaler
	// and BinaryUnmarshaler for state (de)serialization as documented by hash.Hash.
	hash, err := blake2b.New(blake2b.Size256, nil)
	if err != nil {
		panic(err)
	}

	return hash
}

var HashBLAKE2 noise.HashFunc = blake2Fn{}

// ChaCha20-Poly1305 usually offers better performance than the more prevalent AES-GCM algorithm on systems where the CPU(s)
// does not feature the AES-NI instruction set extension.[2] As a result, ChaCha20-Poly1305 is sometimes preferred over
// AES-GCM due to its similar levels of security and in certain use cases involving mobile devices, which mostly use ARM-based CPUs.
var cipherSuite = noise.NewCipherSuite(noise.DH25519, noise.CipherChaChaPoly, HashBLAKE2)

type Noise struct {
}

// handshake
