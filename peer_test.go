package noise

import (
	"bytes"
	"encoding/hex"
	"errors"
	"net"
	"testing"
	"time"
)

// Group of prebuilt peers, public keys and sessions to test purpose
var (
	PeerAPb = PublicKey("f46bea91688c3187eebe66f25f1bcfcb6696c90c293b3a9dca749f6218b7bb52")
	PeerBPb = PublicKey("d0bf26bed4774c612691fd7a618dd23660e316dde3916da5c7698dc9b685e2ae")
	PeerCPb = PublicKey("4c67ad6ef6287f0cf7b1b888c1e93eb4c685e3bc59c33b1ecf79a3ad227219e8")
	PeerDPb = PublicKey("83a2dd209b270d19aedaa4e588fd94fee599b510a49988efd067967ce25053d0")
	PeerEPb = PublicKey("78112677879bb3922a60cbc12ecbc46fdd33e69447df7186f618a0011056a3c1")
	PeerFPb = PublicKey("4c268f42ac66ed02f62d0f8951c7fa042b0a281f57385daf8ee4576b30b8fc00")
)

var (
	sessionA = mockSession(&mockConn{}, PeerAPb)
	sessionB = mockSession(&mockConn{}, PeerBPb)
	sessionC = mockSession(&mockConn{}, PeerCPb)
	sessionD = mockSession(&mockConn{}, PeerDPb)
	sessionE = mockSession(&mockConn{}, PeerEPb)
	sessionF = mockSession(&mockConn{}, PeerFPb)
)

var (
	peerA = newPeer(sessionA)
	peerB = newPeer(sessionB)
	peerC = newPeer(sessionC)
	peerD = newPeer(sessionD)
	peerE = newPeer(sessionE)
	peerF = newPeer(sessionF)
)

// Mock Address from net.Addr
// ref: https://pkg.go.dev/net#Addr
type mockAddr struct {
	addr string
}

func (*mockAddr) Network() string {
	return "tcp"
}

func (m *mockAddr) String() string {
	return m.addr
}

// Mock Address from net.Conn
// ref: https://pkg.go.dev/net#Conn
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

// Mock handshake state for secure session
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

// mockSession create a testable session
func mockSession(conn net.Conn, pb PublicKey) *session {
	return &session{conn, KeyRing{}, pb, nil, nil}
}

// mockID create a new testable id from public key
func mockID(pb PublicKey) ID {
	var id ID
	addr := []byte(pb)
	copy(id[:], addr)
	return id
}

// From content generate a bytes 32 key
func mockBytes(content PublicKey) []byte {
	var expected [32]byte
	copy(expected[:], content)
	return expected[:]
}

func TestByteID(t *testing.T) {
	expected := mockBytes(PeerAPb)
	id := mockID(PeerAPb)

	if !bytes.Equal(id.Bytes(), expected[:]) {
		t.Errorf("expected returned bytes equal to %v, got %v", string(expected[:]), string(id.Bytes()))
	}
}

func TestStringID(t *testing.T) {
	id := mockID(PeerAPb)
	expected := mockBytes(PeerAPb)

	if id.String() != string(expected) {
		t.Errorf("expected returned string equal to %v, got %v", string(PeerAPb), id.String())
	}
}

func TestHashID(t *testing.T) {
	expected := "fab03245b98fc2491b64810d9ab7fccf86db272a54c038780d88852937d25242"
	got := hex.EncodeToString(blake2([]byte(PeerAPb)))

	if expected != got {
		t.Errorf("expected returned %s equal to %s", got, expected)
	}
}
