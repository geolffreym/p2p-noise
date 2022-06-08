package network

import (
	"errors"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/geolffreym/p2p-noise/network"
)

type mockAddr struct{}

func (*mockAddr) Network() string {
	return "tcp"
}

func (*mockAddr) String() string {
	return "127.0.0.1:23"
}

// net/net.go
type mockConn struct {
	channel    chan []byte // Simulation for Message network exchange
	shouldFail bool
}

// Read reads data from the connection.
// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
func (c *mockConn) Read(p []byte) (n int, err error) {
	data := <-c.channel
	return len(data), nil
}

// Write writes data to the connection.
// Write can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetWriteDeadline.
// time limit; see SetDeadline and SetWriteDeadline.
func (c *mockConn) Write(b []byte) (n int, err error) {
	c.channel <- b

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
	return &mockAddr{}
}

// RemoteAddr returns the remote network address, if known.
func (c *mockConn) RemoteAddr() net.Addr {
	return &mockAddr{}
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

func TestSocket(t *testing.T) {
	conn := &mockConn{}
	address := network.Socket("127.0.0.1:23")
	peer := network.NewPeer(address, conn)

	if peer.Socket() != address {
		t.Errorf("expected socket %#v, got %#v", address, peer.Socket())
	}

}

func TestFailClose(t *testing.T) {
	conn := &mockConn{shouldFail: true}
	address := network.Socket("127.0.0.1:23")
	peer := network.NewPeer(address, conn)

	err := peer.Close()

	if err == nil {
		t.Errorf("expected error but got %#v", err)
	}
}

func TestConnection(t *testing.T) {
	conn := &mockConn{}
	address := network.Socket("127.0.0.1:23")
	peer := network.NewPeer(address, conn)

	if !reflect.DeepEqual(peer.Connection(), conn) {
		t.Errorf("expected error but got %#v", peer.Connection())
	}

}

func TestSendReceive(t *testing.T) {

	channel := make(chan []byte)
	conn := &mockConn{channel: channel}

	address := network.Socket("127.0.0.1:23")
	peer := network.NewPeer(address, conn)

	expected := "ping from peer"
	// Someone sending a message to peer
	go func(p network.Peer) {
		p.Send([]byte(expected))
	}(peer)

	// Waiting for incoming messages
	t.Run("Reading", func(t *testing.T) {
		// Simulating network
		buf := make([]byte, 1024)
		bytes, _ := peer.Receive(buf)

		if bytes != len([]byte(expected)) {
			t.Errorf("expected receive same bytes sent \"%s\", got \"%s\"", expected, string(buf))
		}

	})

}
