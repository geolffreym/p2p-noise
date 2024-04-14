package noise

import (
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	var expected []ID
	router := newRouter()
	peers := []*peer{peerA, peerB, peerC, peerD}

	for _, peer := range peers {
		router.Add(peer)
		expected = append(expected, peer.ID())
	}

	// Table driven test
	for _, e := range expected {
		t.Run(fmt.Sprintf("%x", e), func(t *testing.T) {
			// Match recently added peer
			if _, ok := router.Query(e); !ok {
				t.Errorf("expected routed peer id %x", e.String())
			}
		})

	}

}

func TestQuery(t *testing.T) {
	router := newRouter()
	// Add new record
	router.Add(peerA)
	router.Add(peerB)
	expected := []ID{peerA.ID(), peerB.ID()}

	// Table driven test
	for _, e := range expected {
		t.Run(fmt.Sprintf("%x", e), func(t *testing.T) {
			// Return the socket related peer
			if peer, ok := router.Query(e); !ok {
				t.Errorf("expected peer for valid socket %#v, got %v", e.String(), peer)
			}
		})

	}

}

func TestInvalidQuery(t *testing.T) {
	router := newRouter()
	id := mockID(PeerBPb)
	if peer, ok := router.Query(id); ok {
		t.Errorf("expected nil for invalid socket %#v, got %v", PeerBPb, peer)
	}

}

func TestLen(t *testing.T) {
	router := newRouter()
	router.Add(peerA)
	router.Add(peerB)
	router.Add(peerC)
	router.Add(peerD)

	if router.Len() != 4 {
		t.Errorf("expected 4 len for registered peers, got %v", router.Len())
	}

}

func TestTable(t *testing.T) {
	router := newRouter()
	expected := []string{
		peerA.ID().String(),
		peerB.ID().String(),
	}

	router.Add(peerA)
	router.Add(peerB)

LOOP:
	for peer := range router.Table() {
		got := peer.ID().String()

		for _, expect := range expected {
			if expect == got {
				// move to loop and start again with next
				// this approach is equivalent to run a needle in a haystack
				// and avoid the error if match found forwarding the iteration to the main loop
				continue LOOP
			}
		}

		t.Errorf("expected corresponding table matching entry %x", got)
	}

}

func TestDelete(t *testing.T) {
	router := newRouter()

	// Add new record
	router.Add(peerA) // 1
	router.Add(peerB) // 2
	router.Add(peerC) // 3
	router.Add(peerD) // 4
	router.Add(peerE) // 5

	// delete B and F
	router.Remove(peerB)
	router.Remove(peerE)

	if _, ok := router.Query(peerB.ID()); ok {
		t.Errorf("expected %v not registered in router after delete", peerB.ID())
	}

	if _, ok := router.Query(peerE.ID()); ok {
		t.Errorf("expected %v not registered in router after delete", peerE.ID())
	}

}
