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
)

// Default protocol
const PROTOCOL = "tcp"

type NetworkRouter interface {
	Table() Table
}

type NetworkBroker interface {
	Publish(event Event, buf []byte, peer PeerStreamer)
	Register(e Event, s Messenger)
}

type NetworkConnection interface {
	Dial(addr string) error
	Listen(addr string) error
	Close()
}

type NetworkMonitor interface {
	watch(peer Peer)
	Closed() bool
}

type Network interface {
	NetworkRouter
	NetworkBroker
	NetworkConnection
	NetworkMonitor
}

// Network communication logic
// If a type exists only to implement an interface and will never have
// exported methods beyond that interface, there is no need to export the type itself.
// Exporting just the interface makes it clear the value has no interesting behavior
// beyond what is described in the interface.
// It also avoids the need to repeat the documentation on every instance of a common method.
//
// ref: https://go.dev/doc/effective_go#interfaces
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

// watch watchdog for incoming messages.
// incoming message monitor is suggested to be processed in go routines.
func (network *network) watch(peer Peer) {
	buf := make([]byte, 1024)

KEEPALIVE:
	for {
		// Sync buffer reading
		_, err := peer.Receive(buf)
		// If connection is closed
		// stop routines watching peers
		if network.Closed() {
			return
		}

		if err != nil {
			// net: don't return io.EOF from zero byte reads
			// if err == io.EOF then peer connection is closed
			_, isNetError := err.(*net.OpError)
			if err == io.EOF || isNetError {
				err := peer.Close() // Close disconnected peer
				if err != nil {
					log.Fatal(errors.Closing(err).Error())
				}

				//Notify to network about the peer state
				network.Publish(PEER_DISCONNECTED, []byte(peer.Socket()), peer)
				// Remove peer from router table
				network.router.Delete(peer)
				return
			}

			// Keep alive always that zero bytes are not received
			break KEEPALIVE
		}

		// Emit new incoming message notification
		network.Publish(MESSAGE_RECEIVED, buf, peer)
	}

}

// routing initialize route in routing table from connection interface
// Return new peer added to table
func (network *network) routing(conn net.Conn) Peer {

	// Assertion for tcp connection to keep alive
	connection, isTCP := conn.(*net.TCPConn)
	if isTCP {
		// If tcp enforce keep alive connection
		// SetKeepAlive sets whether the operating system should send keep-alive messages on the connection.
		connection.SetKeepAlive(true)
	}

	// Routing connections
	remote := connection.RemoteAddr().String()
	// eg. 192.168.1.1:8080
	socket := Socket(remote)
	// We need to know how interact with peer based on socket and connection
	peer := NewPeer(socket, conn)
	return network.router.Add(peer)
}

// publish emit network event notifications
func (network *network) Publish(event Event, buf []byte, peer PeerStreamer) {
	// Emit new notification
	message := NewMessage(event, buf, peer)
	network.events.Publish(message)
}

// Register associate subscriber to a event channel
// alias for internal Event Register
func (network *network) Register(e Event, s Messenger) {
	network.events.Register(e, s)
}

// Listen start listening on the given address and wait for new connection.
// Return network as nil and error if error occurred while listening.
func (network *network) Listen(addr string) error {
	listener, err := net.Listen(PROTOCOL, addr)
	if err != nil {
		return err
	}

	// Dispatch event on start listening
	network.Publish(SELF_LISTENING, []byte(addr), nil)
	// monitor connection to close listener
	go func(listener net.Listener) {
		<-network.sentinel
		err := listener.Close()
		if err != nil {
			log.Fatal(errors.Closing(err).Error())
		}
	}(listener)

	for {
		// Block/Hold while waiting for new incoming connection
		// Synchronized incoming connections
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(errors.Binding(err).Error())
			return err
		}

		peer := network.routing(conn) // Routing for connection
		go network.watch(peer)        // Wait for incoming messages
		// Dispatch event for new peer connected
		payload := []byte(peer.Socket())
		network.Publish(NEWPEER_DETECTED, payload, peer)
	}

}

// Return current routing table
func (network *network) Table() Table {
	return network.router.Table()
}

// Closed Non-blocking check connection state.
// Return true for connection open else false
func (network *network) Closed() bool {
	select {
	case <-network.sentinel:
		return true
	default:
		return false
	}
}

// Close all peers connections and stop listening
func (network *network) Close() {
	for _, peer := range network.router.Table() {
		go func(p Peer) {
			if err := p.Close(); err != nil {
				log.Fatal(errors.Closing(err).Error())
			}
		}(peer)
	}

	// Dispatch event on close network
	network.Publish(CLOSED_CONNECTION, []byte(""), nil)
	// If channel get closed then all routines waiting for connections
	// or waiting for incoming messages get closed too.
	close(network.sentinel)
}

// Dial to node and add connected peer to routing table
// Return network as nil and error if error occurred while dialing network.
func (network *network) Dial(addr string) error {
	conn, err := net.Dial(PROTOCOL, addr)
	if err != nil {
		return errors.Dialing(err, addr)
	}

	peer := network.routing(conn) // Routing for connection
	go network.watch(peer)        // Wait for incoming messages
	// Dispatch event for new peer connected
	network.Publish(NEWPEER_DETECTED, []byte(peer.Socket()), peer)
	return nil
}
