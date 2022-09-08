package noise

import (
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

func TestAdd(t *testing.T) {
	router := newRouter()
	// Add new record
	router.Add(newPeer(&mockConn{addr: PeerA}))
	router.Add(newPeer(&mockConn{addr: PeerB}))
	router.Add(newPeer(&mockConn{addr: PeerC}))
	router.Add(newPeer(&mockConn{addr: PeerD}))
	router.Add(newPeer(&mockConn{addr: PeerE}))
	router.Add(newPeer(&mockConn{addr: PeerF}))

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
			if _, ok := router.Table()[ID(e.socket)]; !ok {
				t.Errorf("expected routed socket %#v", e.socket)
			}
		})

	}

}

func TestQuery(t *testing.T) {
	router := newRouter()
	// Add new record
	router.Add(newPeer(&mockConn{addr: PeerA}))
	router.Add(newPeer(&mockConn{addr: PeerB}))

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
			id := ID(e.socket)
			if peer := router.Query(id); peer == nil {
				t.Errorf("expected peer for valid socket %#v, got %v", e.socket, peer)
			}
		})

	}

}

func TestInvalidQuery(t *testing.T) {
	router := newRouter()
	// Add new record
	router.Add(newPeer(&mockConn{addr: PeerA}))

	if peer := router.Query(PeerB); peer != nil {
		t.Errorf("expected nil for invalid socket %#v, got %v", PeerB, peer)
	}

}

func TestLen(t *testing.T) {
	router := newRouter()
	// Add new record
	router.Add(newPeer(&mockConn{addr: PeerA})) // 1
	router.Add(newPeer(&mockConn{addr: PeerB})) // 2
	router.Add(newPeer(&mockConn{addr: PeerC})) // 3
	router.Add(newPeer(&mockConn{addr: PeerD})) // 4
	router.Add(newPeer(&mockConn{addr: PeerE})) // 5
	router.Add(newPeer(&mockConn{addr: PeerF})) // 6

	if router.Len() != 6 {
		t.Errorf("expected 6 len for registered peers,  got %v", router.Len())
	}

}

func TestDelete(t *testing.T) {
	router := newRouter()

	peerA := newPeer(&mockConn{addr: PeerA})
	peerB := newPeer(&mockConn{addr: PeerB})
	peerC := newPeer(&mockConn{addr: PeerC})
	peerD := newPeer(&mockConn{addr: PeerD})
	peerE := newPeer(&mockConn{addr: PeerE})

	// Add new record
	router.Add(peerA) // 1
	router.Add(peerB) // 2
	router.Add(peerC) // 3
	router.Add(peerD) // 4
	router.Add(peerE) // 5

	// delete B and F
	router.Remove(peerB)
	router.Remove(peerE)

	if router.Query(peerB.ID()) != nil {
		t.Errorf("expected %v not registered in router after delete", peerB.ID())
	}

	if router.Query(peerE.ID()) != nil {
		t.Errorf("expected %v not registered in router after delete", peerE.ID())
	}

}

func TestFlush(t *testing.T) {
	router := newRouter()
	peerA := newPeer(&mockConn{addr: PeerA})
	// Add new record
	router.Add(peerA) // 1
	router.Flush()

	if router.Table() != nil {
		t.Errorf("expected empty table, got %v", router.Table())
	}

}

func TestFlushSize(t *testing.T) {
	router := newRouter()

	peerA := newPeer(&mockConn{addr: PeerA})
	peerB := newPeer(&mockConn{addr: PeerB})
	peerC := newPeer(&mockConn{addr: PeerC})

	// Add new record
	router.Add(peerA) // 1
	router.Add(peerB) // 2
	router.Add(peerC) // 3

	// delete B and F
	len := router.Len()
	deleted := router.Flush()

	if deleted != len {
		t.Errorf("expected %v table flushed peers, got %v", len, deleted)
	}

}
