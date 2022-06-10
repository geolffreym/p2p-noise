package noise

import (
	"testing"
)

func TestAdd(t *testing.T) {
	router := newRouter()
	// Add new record
	router.Add(newPeer("127.0.0.1:8080", nil))
	router.Add(newPeer("127.0.0.1:8081", nil))
	router.Add(newPeer("127.0.0.1:8082", nil))
	router.Add(newPeer("127.0.0.1:8083", nil))
	router.Add(newPeer("127.0.0.1:8084", nil))
	router.Add(newPeer("127.0.0.1:8085", nil))

	expected := []struct {
		socket string
	}{
		{
			socket: "127.0.0.1:8080",
		},
		{
			socket: "127.0.0.1:8081",
		},
		{
			socket: "127.0.0.1:8082",
		},
		{
			socket: "127.0.0.1:8083",
		},
		{
			socket: "127.0.0.1:8084",
		},
		{
			socket: "127.0.0.1:8085",
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
	router.Add(newPeer("127.0.0.1:8080", nil))
	router.Add(newPeer("127.0.0.1:8081", nil))

	expected := []struct {
		socket string
	}{
		{
			socket: "127.0.0.1:8080",
		},
		{
			socket: "127.0.0.1:8081",
		},
	}

	// Table driven test
	// For each expected event
	for _, e := range expected {
		t.Run(e.socket, func(t *testing.T) {

			if peer := router.Query(Socket(e.socket)); peer == nil {
				t.Errorf("expected peer for valid socket %#v, , got %sv", e.socket, peer)
			}
		})

	}

}

func TestInvalidQuery(t *testing.T) {
	router := newRouter()
	// Add new record
	router.Add(newPeer("127.0.0.1:8080", nil))

	if peer := router.Query(Socket("127.0.0.1:8081")); peer != nil {
		t.Errorf("expected nil for invalid socket %#v, got %sv", "127.0.0.1:8081", peer)
	}

}

func TestLen(t *testing.T) {
	router := newRouter()
	// Add new record
	router.Add(newPeer("127.0.0.1:8080", nil)) // 1
	router.Add(newPeer("127.0.0.1:8081", nil)) // 2
	router.Add(newPeer("127.0.0.1:8082", nil)) // 3
	router.Add(newPeer("127.0.0.1:8083", nil)) // 4
	router.Add(newPeer("127.0.0.1:8084", nil)) // 5
	router.Add(newPeer("127.0.0.1:8085", nil)) // 6

	if router.Len() != 6 {
		t.Errorf("expected 6 len for registered peers,  got %v", router.Len())
	}

}

func TestDelete(t *testing.T) {
	router := newRouter()

	peerA := newPeer("127.0.0.1:8080", nil)
	peerB := newPeer("127.0.0.1:8080", nil)
	peerC := newPeer("127.0.0.1:8080", nil)
	peerD := newPeer("127.0.0.1:8080", nil)
	peerF := newPeer("127.0.0.1:8080", nil)

	// Add new record
	router.Add(peerA) // 1
	router.Add(peerB) // 2
	router.Add(peerC) // 3
	router.Add(peerD) // 4
	router.Add(peerF) // 5

	// delete B and F
	router.Delete(peerB)
	router.Delete(peerF)

	if router.Query(peerB.Socket()) != nil {
		t.Errorf("expected %v not registered in router after delete", peerB.Socket())
	}

	if router.Query(peerF.Socket()) != nil {
		t.Errorf("expected %v not registered in router after delete", peerF.Socket())
	}

}
