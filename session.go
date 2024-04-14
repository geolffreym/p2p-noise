package noise

import (
	"crypto/ed25519"
	"net"

	"golang.org/x/crypto/blake2b"
)

// blake2 return a 32-bytes 256 checksum representation for blake2 hash.
func blake2(i []byte) []byte {
	// New returns a new hash.Hash computing the BLAKE2b checksum with a custom length.
	// A non-nil key turns the hash into a MAC. The key must be between zero and 64 bytes long.
	// The hash size can be a value between 1 and 64 but it is highly recommended to use
	// values equal or greater than:
	// - 32 if BLAKE2b is used as a hash function (The key is zero bytes long).
	// - 16 if BLAKE2b is used as a MAC function (The key is at least 16 bytes long).
	// When the key is nil, the returned hash.Hash implements BinaryMarshaler
	// and BinaryUnmarshaler for state (de)serialization as documented by hash.Hash.
	hash, err := blake2b.New256(nil)
	if err != nil {
		return nil
	}

	hash.Write(i)
	return hash.Sum(nil)
}

// session implement net.Conn interface.
// session methods to provide secure message exchange between peers.
// Please see [Connection Interface] for more details.
//
// [Connection Interface]: https://pkg.go.dev/net#Conn
type session struct {
	net.Conn             // insecure conn
	kr         KeyRing   // local bundle of keys
	svk        PublicKey // remote public key
	encryption CipherState
	decryption CipherState
}

// Create a new secure session
func newSession(conn net.Conn, kr KeyRing) (*session, error) {
	return &session{conn, kr, PublicKey{}, nil, nil}, nil
}

// Set encryption/decryption state for session.
// A CipherState provides symmetric encryption and decryption after a successful handshake
func (s *session) SetCyphers(enc, dec CipherState) {
	s.encryption = enc // pb-k
	s.decryption = dec // pv-k
}

// SetVerifyKey set remote signature validation public key.
func (s *session) SetRemotePublicKey(pb PublicKey) {
	s.svk = pb
}

// Encrypt cipher message using encryption keys provided in handshake.
func (s *session) Encrypt(out, msg []byte) ([]byte, error) {
	return s.encryption.Encrypt(out, nil, msg)
}

// Decrypt message using decryption keys provided in handshake.
func (s *session) Decrypt(out, digest []byte) ([]byte, error) {
	return s.decryption.Decrypt(out, nil, digest)
}

// RemotePublicKey returns the static key provided by the remote peer during a handshake.
func (s *session) RemotePublicKey() []byte {
	return s.svk
}

// Sign message with local private key.
func (s *session) Sign(msg []byte) []byte {
	return ed25519.Sign(s.kr.sv.Private, msg)
}

// Verify message with remote public key.
func (s *session) Verify(msg, sig []byte) bool {
	// Use remote peer public key to verify message
	return ed25519.Verify(s.svk, msg, sig)
}
