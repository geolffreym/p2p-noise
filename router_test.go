package noise

import (
	"bytes"
	"testing"
)

const (
	PeerA = "127.0.0.1:8080"
	PeerB = "127.0.0.1:8081"
	PeerC = "127.0.0.1:8082"
	PeerD = "127.0.0.1:8083"
	PeerE = "127.0.0.1:8084"
	PeerF = "127.0.0.1:8085"
)

func TestByteSocket(t *testing.T) {
	socket := Socket(LOCAL_ADDRESS)
	if !bytes.Equal(socket.Bytes(), []byte(LOCAL_ADDRESS)) {
		t.Errorf("Expected returned bytes equal to %v", string(socket.Bytes()))
	}
}

func TestStringSocket(t *testing.T) {
	socket := Socket(LOCAL_ADDRESS)
	if socket.String() != LOCAL_ADDRESS {
		t.Errorf("Expected returned string equal to %v", socket.String())
	}
}

func TestAdd(t *testing.T) {
	router := newRouter()
	// Add new record
	router.Add(newPeer(PeerA, &mockConn{}))
	router.Add(newPeer(PeerB, &mockConn{}))
	router.Add(newPeer(PeerC, &mockConn{}))
	router.Add(newPeer(PeerD, &mockConn{}))
	router.Add(newPeer(PeerE, &mockConn{}))
	router.Add(newPeer(PeerF, &mockConn{}))

	expected := []struct {
		socket string
	}{
		{
			socket: PeerA,
		},
		{
			socket: PeerB,
		},
		{
			socket: PeerC,
		},
		{
			socket: PeerD,
		},
		{
			socket: PeerE,
		},
		{
			socket: PeerF,
		},
	}

	// Table driven test
	for _, e := range expected {
		t.Run(e.socket, func(t *testing.T) {
			// Match recently added peer
			if _, ok := router.Table()[Socket(e.socket)]; !ok {
				t.Errorf("expected routed socket %#v", e.socket)
			}
		})

	}

}

func TestQuery(t *testing.T) {
	router := newRouter()
	// Add new record
	router.Add(newPeer(PeerA, &mockConn{}))
	router.Add(newPeer(PeerB, &mockConn{}))

	expected := []struct {
		socket string
	}{
		{
			socket: PeerA,
		},
		{
			socket: PeerB,
		},
	}

	// Table driven test
	for _, e := range expected {
		t.Run(e.socket, func(t *testing.T) {
			// Return the socket related peer
			if peer := router.Query(Socket(e.socket)); peer == nil {
				t.Errorf("expected peer for valid socket %#v, got %v", e.socket, peer)
			}
		})

	}

}

func TestInvalidQuery(t *testing.T) {
	router := newRouter()
	// Add new record
	router.Add(newPeer(PeerA, &mockConn{}))

	if peer := router.Query(Socket(PeerB)); peer != nil {
		t.Errorf("expected nil for invalid socket %#v, got %v", PeerB, peer)
	}

}

func TestLen(t *testing.T) {
	router := newRouter()
	// Add new record
	router.Add(newPeer(PeerA, &mockConn{})) // 1
	router.Add(newPeer(PeerB, &mockConn{})) // 2
	router.Add(newPeer(PeerC, &mockConn{})) // 3
	router.Add(newPeer(PeerD, &mockConn{})) // 4
	router.Add(newPeer(PeerE, &mockConn{})) // 5
	router.Add(newPeer(PeerF, &mockConn{})) // 6

	if router.Len() != 6 {
		t.Errorf("expected 6 len for registered peers,  got %v", router.Len())
	}

}

func TestDelete(t *testing.T) {
	router := newRouter()

	peerA := newPeer(PeerA, &mockConn{})
	peerB := newPeer(PeerB, &mockConn{})
	peerC := newPeer(PeerC, &mockConn{})
	peerD := newPeer(PeerD, &mockConn{})
	peerE := newPeer(PeerE, &mockConn{})

	// Add new record
	router.Add(peerA) // 1
	router.Add(peerB) // 2
	router.Add(peerC) // 3
	router.Add(peerD) // 4
	router.Add(peerE) // 5

	// delete B and F
	router.Remove(peerB)
	router.Remove(peerE)

	if router.Query(peerB.Socket()) != nil {
		t.Errorf("expected %v not registered in router after delete", peerB.Socket())
	}

	if router.Query(peerE.Socket()) != nil {
		t.Errorf("expected %v not registered in router after delete", peerE.Socket())
	}

}

func TestFlush(t *testing.T) {
	router := newRouter()
	peerA := newPeer(PeerA, &mockConn{})
	// Add new record
	router.Add(peerA) // 1
	router.Flush()

	if router.Table() != nil {
		t.Errorf("expected empty table, got %v", router.Table())
	}

}

func TestFlushSize(t *testing.T) {
	router := newRouter()

	peerA := newPeer(PeerA, &mockConn{})
	peerB := newPeer(PeerB, &mockConn{})
	peerC := newPeer(PeerC, &mockConn{})
	peerD := newPeer(PeerD, &mockConn{})
	peerE := newPeer(PeerE, &mockConn{})

	// Add new record
	router.Add(peerA) // 1
	router.Add(peerB) // 2
	router.Add(peerC) // 3
	router.Add(peerD) // 4
	router.Add(peerE) // 5

	// delete B and F
	len := router.Len()
	deleted := router.Flush()

	if deleted != len {
		t.Errorf("expected %v table flushed peers, got %v", len, deleted)
	}

}
