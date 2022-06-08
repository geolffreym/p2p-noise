package network_test

import (
	"testing"

	"github.com/geolffreym/p2p-noise/network"
)

func TestAdd(t *testing.T) {
	router := network.NewRouter()
	// Add new record
	router.Add(network.NewPeer("127.0.0.1:8080", nil))
	router.Add(network.NewPeer("127.0.0.1:8081", nil))
	router.Add(network.NewPeer("127.0.0.1:8082", nil))
	router.Add(network.NewPeer("127.0.0.1:8083", nil))
	router.Add(network.NewPeer("127.0.0.1:8084", nil))
	router.Add(network.NewPeer("127.0.0.1:8085", nil))

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

			if _, ok := router.Table()[network.Socket(e.socket)]; !ok {
				t.Errorf("expected routed socket %#v", e.socket)
			}
		})

	}

}

func TestQuery(t *testing.T) {
	router := network.NewRouter()
	// Add new record
	router.Add(network.NewPeer("127.0.0.1:8080", nil))
	router.Add(network.NewPeer("127.0.0.1:8081", nil))

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

			if peer := router.Query(network.Socket(e.socket)); peer == nil {
				t.Errorf("expected peer for valid socket %#v, , got %sv", e.socket, peer)
			}
		})

	}

}

func TestInvalidQuery(t *testing.T) {
	router := network.NewRouter()
	// Add new record
	router.Add(network.NewPeer("127.0.0.1:8080", nil))

	if peer := router.Query(network.Socket("127.0.0.1:8081")); peer != nil {
		t.Errorf("expected nil for invalid socket %#v, got %sv", "127.0.0.1:8081", peer)
	}

}

func TestLen(t *testing.T) {
	router := network.NewRouter()
	// Add new record
	router.Add(network.NewPeer("127.0.0.1:8080", nil)) // 1
	router.Add(network.NewPeer("127.0.0.1:8081", nil)) // 2
	router.Add(network.NewPeer("127.0.0.1:8082", nil)) // 3
	router.Add(network.NewPeer("127.0.0.1:8083", nil)) // 4
	router.Add(network.NewPeer("127.0.0.1:8084", nil)) // 5
	router.Add(network.NewPeer("127.0.0.1:8085", nil)) // 6

	if router.Len() != 6 {
		t.Errorf("expected 6 len for registered peers,  got %v", router.Len())
	}

}

func TestDelete(t *testing.T) {
	router := network.NewRouter()

	peerA := network.NewPeer("127.0.0.1:8080", nil)
	peerB := network.NewPeer("127.0.0.1:8080", nil)
	peerC := network.NewPeer("127.0.0.1:8080", nil)
	peerD := network.NewPeer("127.0.0.1:8080", nil)
	peerF := network.NewPeer("127.0.0.1:8080", nil)

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
