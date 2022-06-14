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
	router.Add(newPeer(PeerA, nil))
	router.Add(newPeer(PeerB, nil))
	router.Add(newPeer(PeerC, nil))
	router.Add(newPeer(PeerD, nil))
	router.Add(newPeer(PeerE, nil))
	router.Add(newPeer(PeerF, nil))

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
	// For each expected event
	for _, e := range expected {
		t.Run(e.socket, func(t *testing.T) {

			if _, ok := router.Table()[Socket(e.socket)]; !ok {
				t.Errorf("expected routed socket %#v", e.socket)
			}
		})

	}

}

func TestQuery(t *testing.T) {
	router := newRouter()
	// Add new record
	router.Add(newPeer(PeerA, nil))
	router.Add(newPeer(PeerB, nil))

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
	// For each expected event
	for _, e := range expected {
		t.Run(e.socket, func(t *testing.T) {

			if peer := router.Query(Socket(e.socket)); peer == nil {
				t.Errorf("expected peer for valid socket %#v, got %v", e.socket, peer)
			}
		})

	}

}

func TestInvalidQuery(t *testing.T) {
	router := newRouter()
	// Add new record
	router.Add(newPeer(PeerA, nil))

	if peer := router.Query(Socket(PeerB)); peer != nil {
		t.Errorf("expected nil for invalid socket %#v, got %v", PeerB, peer)
	}

}

func TestLen(t *testing.T) {
	router := newRouter()
	// Add new record
	router.Add(newPeer(PeerA, nil)) // 1
	router.Add(newPeer(PeerB, nil)) // 2
	router.Add(newPeer(PeerC, nil)) // 3
	router.Add(newPeer(PeerD, nil)) // 4
	router.Add(newPeer(PeerE, nil)) // 5
	router.Add(newPeer(PeerF, nil)) // 6

	if router.Len() != 6 {
		t.Errorf("expected 6 len for registered peers,  got %v", router.Len())
	}

}

func TestDelete(t *testing.T) {
	router := newRouter()

	peerA := newPeer(PeerA, nil)
	peerB := newPeer(PeerA, nil)
	peerC := newPeer(PeerA, nil)
	peerD := newPeer(PeerA, nil)
	peerF := newPeer(PeerA, nil)

	// Add new record
	router.Add(peerA) // 1
	router.Add(peerB) // 2
	router.Add(peerC) // 3
	router.Add(peerD) // 4
	router.Add(peerF) // 5

	// delete B and F
	router.Remove(peerB)
	router.Remove(peerF)

	if router.Query(peerB.Socket()) != nil {
		t.Errorf("expected %v not registered in router after delete", peerB.Socket())
	}

	if router.Query(peerF.Socket()) != nil {
		t.Errorf("expected %v not registered in router after delete", peerF.Socket())
	}

}
