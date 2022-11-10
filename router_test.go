package noise

import (
	"fmt"
	"testing"
)

const (
	PeerA = "f46bea91688c3187eebe66f25f1bcfcb6696c90c293b3a9dca749f6218b7bb52"
	PeerB = "d0bf26bed4774c612691fd7a618dd23660e316dde3916da5c7698dc9b685e2ae"
	PeerC = "4c67ad6ef6287f0cf7b1b888c1e93eb4c685e3bc59c33b1ecf79a3ad227219e8"
	PeerD = "83a2dd209b270d19aedaa4e588fd94fee599b510a49988efd067967ce25053d0"
	PeerE = "78112677879bb3922a60cbc12ecbc46fdd33e69447df7186f618a0011056a3c1"
	PeerF = "4c268f42ac66ed02f62d0f8951c7fa042b0a281f57385daf8ee4576b30b8fc00"
)

var (
	sessionA = mockSession(&mockConn{addr: PeerA})
	sessionB = mockSession(&mockConn{addr: PeerB})
	sessionC = mockSession(&mockConn{addr: PeerC})
	sessionD = mockSession(&mockConn{addr: PeerD})
	sessionE = mockSession(&mockConn{addr: PeerE})
	sessionF = mockSession(&mockConn{addr: PeerF})
)

var (
	peerA = newPeer(sessionA)
	peerB = newPeer(sessionB)
	peerC = newPeer(sessionC)
	peerD = newPeer(sessionD)
	peerE = newPeer(sessionE)
	peerF = newPeer(sessionF)
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
			if p := router.Query(e); p != nil {
				t.Errorf("expected routed socket %#v", e.String())
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
			if peer := router.Query(e); peer == nil {
				t.Errorf("expected peer for valid socket %#v, got %v", e.String(), peer)
			}
		})

	}

}

func TestInvalidQuery(t *testing.T) {
	router := newRouter()
	id := mockID(PeerB)
	if peer := router.Query(id); peer != nil {
		t.Errorf("expected nil for invalid socket %#v, got %v", PeerB, peer)
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

func TestDelete(t *testing.T) {
	router := newRouter()

	peerA := newPeer(mockSession(&mockConn{addr: PeerA}))
	peerB := newPeer(mockSession(&mockConn{addr: PeerB}))
	peerC := newPeer(mockSession(&mockConn{addr: PeerC}))
	peerD := newPeer(mockSession(&mockConn{addr: PeerD}))
	peerE := newPeer(mockSession(&mockConn{addr: PeerE}))

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
