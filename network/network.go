// Package network implements a lightweight TCP communication.
// Offers pretty basic features to communicate between nodes.
//
// See also: https://pkg.go.dev/net#Conn
package network

import (
	"io"
	"log"
	"net"

	"github.com/geolffreym/p2p-noise/errors"
	"github.com/geolffreym/p2p-noise/pubsub"
)

// Default protocol
const PROTOCOL = "tcp"

// Network communication logic
type Network struct {
	table    Router         // Routing hash table eg. {Socket: Conn interface}.
	sentinel chan bool      // Channel flag waiting for signal to close connection.
	Events   pubsub.Channel // Pubsub notifications.
}

// Network factory.
func New() *Network {
	return &Network{
		table:  make(Router),
		Events: make(pubsub.Channel),
	}
}

// Build a new peer from network connection
func (network *Network) peer(conn net.Conn) *Peer {
	return &Peer{
		conn:   conn,
		socket: Socket(conn.RemoteAddr().String()),
	}
}

// Initialize route in routing table from connection interface
// Return new peer added to table
func (network *Network) routing(conn net.Conn) *Peer {
	// Keep routing for each connection
	socket := Socket(conn.RemoteAddr().String())
	peer := network.peer(conn)
	network.table.Add(socket, peer)
	return peer
}

// Run routed stream message in goroutine.
// Each incoming message is processed in non-blocking approach.
func (network *Network) stream(peer *Peer) {
	go func(n *Network, p *Peer) {
		buf := make([]byte, 1024)

	KEEPALIVE:
		for {
			// Stop routine
			if n.IsClosed() {
				return
			}

			_, err := p.Read(buf)
			if err != nil {
				if err == io.EOF {
					break KEEPALIVE
				}
			}

			// TODO Need refactor to handle biggest messages
			// Emit new incoming
			message := pubsub.NewMessage(pubsub.MESSAGE_RECEIVED, buf)
			n.Events.Publish(message)

		}
	}(network, peer)
}

// Concurrent `Bind` network for streams.
// Start a new goroutine to keep waiting for new connections.
func (network *Network) bind(listener net.Listener) {
	go func(n *Network, l net.Listener) {
		for {
			// Block/Hold while waiting for new incoming connection
			// Synchronized incoming connections
			conn, err := l.Accept()
			if err != nil || n.IsClosed() {
				log.Fatalf(errors.WrapBinding(err).Error())
				return
			}

			// Routing for connection
			peer := n.routing(conn)
			n.stream(peer)
			// Dispatch event
			payload := []byte(peer.Socket())
			message := pubsub.NewMessage(pubsub.NEWPEER_DETECTED, payload)
			n.Events.Publish(message)
		}
	}(network, listener)
}

// Start listening on the given address and wait for new connection.
// Return network as nil and error if error occurred while listening.
func (network *Network) Listen(addr string) (*Network, error) {
	listener, err := net.Listen(PROTOCOL, addr)
	if err != nil {
		return nil, errors.WrapListen(err, addr)
	}

	// Concurrent processing for each incoming connection
	network.bind(listener)
	// Dispatch event on start listening
	payload := []byte(addr)
	message := pubsub.NewMessage(pubsub.SELF_LISTENING, payload)
	network.Events.Publish(message)
	return network, nil
}

// Return current routing table
func (network *Network) Table() Router {
	return network.table
}

// Non-blocking check connection state.
// true for connection open else close
func (network *Network) IsClosed() bool {
	select {
	case <-network.sentinel:
		return true
	default:
	}

	return false
}

// Close all peers connections
func (network *Network) Close() {
	for _, peer := range network.table {
		go func(p *Peer) {
			if err := p.Close(); err != nil {
				log.Fatalf(errors.WrapClose(err).Error())
			}
		}(peer)
	}

	// If channel get closed then all routines waiting for connections
	// or waiting for incoming messages get closed too.
	close(network.sentinel)
}

// Dial to a network node and add peer to table
// Return network as nil and error if error occurred while dialing network.
func (network *Network) Dial(addr string) (*Network, error) {
	conn, err := net.Dial(PROTOCOL, addr)
	if err != nil {
		return nil, errors.WrapDial(err, addr)
	}

	// Routing for connection
	peer := network.routing(conn)
	network.stream(peer)
	// Dispatch event
	payload := []byte(peer.Socket())
	message := pubsub.NewMessage(pubsub.NEWPEER_DETECTED, payload)
	network.Events.Publish(message)
	return network, nil
}
