package noise

import (
	"bytes"
	"errors"
	"net"
	"testing"
	"time"
)

const LOCAL_ADDRESS = "127.0.0.1:23"

type mockAddr struct {
	addr string
}

func (*mockAddr) Network() string {
	return "tcp"
}

func (m *mockAddr) String() string {
	return m.addr
}

// net/net.go
type mockConn struct {
	addr       string
	shouldFail bool
	msg        []byte
}

// Read reads data from the connection.
// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
func (c *mockConn) Read(p []byte) (n int, err error) {
	return 0, nil
}

// Write writes data to the connection.
// Write can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetWriteDeadline.
// time limit; see SetDeadline and SetWriteDeadline.
func (c *mockConn) Write(b []byte) (n int, err error) {
	return len(c.msg), nil
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (c *mockConn) Close() error {
	if c.shouldFail {
		return errors.New("failing")
	}
	return nil
}

// LocalAddr returns the local network address, if known.
func (c *mockConn) LocalAddr() net.Addr {
	return &mockAddr{c.addr}
}

// RemoteAddr returns the remote network address, if known.
func (c *mockConn) RemoteAddr() net.Addr {
	return &mockAddr{c.addr}
}

// A zero value for t means I/O operations will not time out.
func (c *mockConn) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (c *mockConn) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (c *mockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

type mockHandshakeState struct {
	addr string
}

func (m *mockHandshakeState) PeerStatic() []byte {
	return []byte(m.addr)
}

func (*mockHandshakeState) MessageIndex() int {
	return 0
}

func (*mockHandshakeState) WriteMessage(out, payload []byte) ([]byte, CipherState, CipherState, error) {
	return nil, nil, nil, nil
}

func (*mockHandshakeState) ReadMessage(out, message []byte) ([]byte, CipherState, CipherState, error) {
	return nil, nil, nil, nil
}

func mockSession(conn net.Conn) *session {
	return &session{
		conn, nil, nil,
		&mockHandshakeState{conn.RemoteAddr().String()},
	}
}

func mockID(address string) ID {
	var id ID
	addr := []byte(address)
	copy(id[:], addr)
	return id
}

func mockBytes(content string) []byte {
	var expected [32]byte
	copy(expected[:], content)
	return expected[:]
}

func TestByteID(t *testing.T) {
	expected := mockBytes(LOCAL_ADDRESS)
	id := mockID(LOCAL_ADDRESS)

	if !bytes.Equal(id.Bytes(), expected[:]) {
		t.Errorf("expected returned bytes equal to %v, got %v", string(expected[:]), string(id.Bytes()))
	}
}

func TestStringID(t *testing.T) {
	id := mockID(LOCAL_ADDRESS)
	expected := mockBytes(LOCAL_ADDRESS)

	if id.String() != string(expected) {
		t.Errorf("expected returned string equal to %v, got %v", string(LOCAL_ADDRESS), id.String())
	}
}
