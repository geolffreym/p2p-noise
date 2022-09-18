package noise

import (
	"net"
)

// session implement net.Conn interface.
// session methods to provide secure message exchange between peers.
// Please see [Connection Interface] for more details.
//
// [Connection Interface]: https://pkg.go.dev/net#Conn
type session struct {
	net.Conn   // insecure conn
	encryption CipherState
	decryption CipherState
	hs         HandshakeState
}

func newSession(conn net.Conn) *session {
	return &session{conn, nil, nil, nil}
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
