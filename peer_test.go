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
	return 1, nil
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

func TestByteID(t *testing.T) {
	id := ID(LOCAL_ADDRESS)
	if !bytes.Equal(id.Bytes(), []byte(LOCAL_ADDRESS)) {
		t.Errorf("Expected returned bytes equal to %v", string(id.Bytes()))
	}
}

func TestStringID(t *testing.T) {
	id := ID(LOCAL_ADDRESS)
	if id.String() != LOCAL_ADDRESS {
		t.Errorf("Expected returned string equal to %v", id.String())
	}
}

func TestID(t *testing.T) {
	address := LOCAL_ADDRESS
	conn := &mockConn{addr: address}
	peer := newPeer(conn)

	if peer.ID() != ID(address) {
		t.Errorf("expected socket %#v, got %#v", address, peer.ID())
	}
}
