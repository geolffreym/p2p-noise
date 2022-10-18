package noise

import (
	"crypto/ed25519"
	"net"

	"golang.org/x/crypto/blake2b"
)

// blake2 return a 32-bytes representation for blake2 hash.
func blake2(i []byte) []byte {
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
	net.Conn         // insecure conn
	kp         DHKey // Diffie-Hellman key pair
	encryption CipherState
	decryption CipherState
	hs         HandshakeState
}

// Create a new secure session
func newSession(conn net.Conn, kp DHKey) *session {
	return &session{conn, kp, nil, nil, nil}
}

// Set encryption/decryption state for session.
// A CipherState provides symmetric encryption and decryption after a successful handshake
func (s *session) SetCyphers(enc, dec CipherState) {
	s.encryption = enc // pb-k
	s.decryption = dec // pv-k
}

// SetState add handshake state to session.
// The state hold all the keys and information about handshake process.
func (s *session) SetState(state HandshakeState) {
	s.hs = state
}

// State return session handshake state
func (s *session) State() HandshakeState {
	return s.hs
}

// Encrypt cipher message using encryption keys provided in handshake.
func (s *session) Encrypt(msg []byte) ([]byte, error) {
	return s.encryption.Encrypt(msg, nil, nil)
}

// Decrypt message using decryption keys provided in handshake.
func (s *session) Decrypt(digest []byte) ([]byte, error) {
	return s.decryption.Decrypt(digest, nil, nil)
}

// Sign message with private key
func (s *session) Sign(msg []byte) []byte {
	return ed25519.Sign(s.kp.Private, msg)
}

// Verify message with remote public key
func (s *session) Verify(msg []byte, sig []byte) bool {
	// Use remote peer public key to verify message
	return ed25519.Verify(s.hs.PeerStatic(), msg, sig)
}
