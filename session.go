package noise

import (
	"net"

	"github.com/flynn/noise"
)

// session implement net.Conn interface.
// session methods to provide secure message exchange between peers.
// Please see [Connection Interface] for more details.
//
// [Connection Interface]: https://pkg.go.dev/net#Conn
type session struct {
	net.Conn   // insecure conn
	encryption *noise.CipherState
	decryption *noise.CipherState
}

func newSession(conn net.Conn) *session {
	return &session{conn, nil, nil}
}

// Set encryption/decryption state for session.
// A CipherState provides symmetric encryption and decryption after a successful handshake
func (s *session) SetCyphers(enc, dec *noise.CipherState) {
	s.encryption = enc // pb-k
	s.decryption = dec // pv-k
}
