// Network implements a lightweight TCP communication.
// Offers pretty basic features to communicate between nodes.
//
// See also: https://pkg.go.dev/net#Conn
package network

import (
	"io"
	"log"
	"net"
	"sync"

	"github.com/geolffreym/p2p-noise/errors"
	"github.com/geolffreym/p2p-noise/utils"
)

// Default protocol
const PROTOCOL = "tcp"

type NetworkRouter interface {
	routing(conn net.Conn) Peer
	Router() Router
}

type NetworkBroker interface {
	observe(peer Peer)
	Publish(event Event, buf []byte, peer PeerStreamer)
	Register(e Event, s Subscriber)
}

type NetworkConnection interface {
	Listen(addr string) (Network, error)
	Dial(addr string) (Network, error)
	Close()
}

type NetworkMonitor interface {
	Closed() bool
}

type Network interface {
	NetworkRouter
	NetworkBroker
	NetworkConnection
	NetworkMonitor
}

// Network communication logic
type network struct {
	sync.RWMutex
	sentinel chan bool // Channel flag waiting for signal to close connection.
	router   Router    // Routing hash table eg. {Socket: Conn interface}.
	events   Events    // Pubsub notifications.
}

// Network factory.
func New() Network {
	return &network{
		sentinel: make(chan bool),
		router:   NewRouter(),
		events:   NewEvents(),
	}
}

// routing initialize route in routing table from connection interface
// Return new peer added to table
func (net *network) routing(conn net.Conn) Peer {

	// Keep routing for each connection
	socket := Socket(conn.RemoteAddr().String())
	peer := NewPeer(socket, conn)
	net.router.Add(socket, peer)
	return peer
}

// publish emit network event notifications
func (network *network) Publish(event Event, buf []byte, peer PeerStreamer) {
	// Emit new notification
	message := NewMessage(event, buf, peer)
	network.events.Publish(message)
}

//  Register associate subscriber to a event channel
//  alias for internal Event Register
func (network *network) Register(e Event, s Subscriber) {
	network.events.Register(e, s)
}

// observe run goroutine waiting for incoming messages.
// Each incoming message is processed in non-blocking approach.
func (network *network) observe(peer Peer) {
	go func(n Network, p Peer) {
		buf := make([]byte, 1024)
	KEEPALIVE:
		for {
			// Stop routine
			if n.Closed() {
				return
			}

			// Sync buffer reading
			_, err := p.Receive(buf)
			if err != nil {
				if err != io.EOF {
					break KEEPALIVE
				}
			}

			// Emit new incoming message notification
			n.Publish(MESSAGE_RECEIVED, buf, p)

		}
	}(network, peer)
}

// bind concurrent network for streams.
// Start a new goroutine to keep waiting for new connections.
func (network *network) bind(listener net.Listener) {
	go func(n Network, l net.Listener) {
		for {
			// Block/Hold while waiting for new incoming connection
			// Synchronized incoming connections
			conn, err := l.Accept()
			if err != nil || n.Closed() {
				log.Fatalf(errors.Binding(err).Error())
				return
			}

			peer := n.routing(conn) // Routing for connection
			n.observe(peer)         // Wait for incoming messages

			// Dispatch event for new peer connected
			payload := []byte(peer.Socket())
			n.Publish(NEWPEER_DETECTED, payload, peer)
		}
	}(network, listener)
}

// Listen start listening on the given address and wait for new connection.
// Return network as nil and error if error occurred while listening.
func (network *network) Listen(addr string) (Network, error) {
	listener, err := net.Listen(PROTOCOL, addr)
	if err != nil {
		return nil, errors.Listening(err, addr)
	}

	network.bind(listener) // Wait for incoming messages
	// Dispatch event on start listening
	network.Publish(SELF_LISTENING, []byte(addr), nil)
	return network, nil
}

// Return current routing table
func (net *network) Router() Router {
	return net.router
}

// Closed Non-blocking check connection state.
// true for connection open else close
func (net *network) Closed() bool {
	select {
	case <-net.sentinel:
		return true
	default:
	}

	return false
}

// Close all peers connections and destroy current state
func (network *network) Close() {
	for _, peer := range network.router.Table() {
		go func(p Peer) {
			if err := p.Close(); err != nil {
				log.Fatalf(errors.Closing(err).Error())
			}
		}(peer)
	}

	// Clear current state after closed connections
	utils.Clear(&network.router)
	utils.Clear(&network.events)
	// Dispatch event on close network
	network.Publish(CLOSED_CONNECTION, []byte(""), nil)
	// If channel get closed then all routines waiting for connections
	// or waiting for incoming messages get closed too.
	close(network.sentinel)
}

// Dial to node and add connected peer to routing table
// Return network as nil and error if error occurred while dialing network.
func (network *network) Dial(addr string) (Network, error) {
	conn, err := net.Dial(PROTOCOL, addr)
	if err != nil {
		return nil, errors.Dialing(err, addr)
	}

	peer := network.routing(conn) // Routing for connection
	network.observe(peer)         // Wait for incoming messages

	// Dispatch event for new peer connected
	network.Publish(NEWPEER_DETECTED, []byte(peer.Socket()), peer)
	return network, nil
}
