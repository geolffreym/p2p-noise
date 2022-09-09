//Copyright (c) 2022, Geolffrey Mena <gmjun2000@gmail.com>

// P2P Noise Library.
// Please read more about [Noise Protocol].
//
// [Noise Protocol]: http://www.noiseprotocol.org/noise.html
package noise

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"
)

// futureDeadline calculate a new time for deadline since now.
func futureDeadLine(deadline time.Duration) time.Time {
	return time.Now().Add(deadline * time.Second)
}

type Config interface {
	// Default "tcp"
	Protocol() string
	// Default 0.0.0.0:8010
	SelfListeningAddress() string
	// Default 100
	MaxPeersConnected() uint8
	// Default 10 << 20 = 10MB
	MaxPayloadSize() uint32
	// Default 1800 seconds = 30 minutes
	PeerDeadline() time.Duration
	// Default 5 seconds
	DialTimeout() time.Duration
}

type Node struct {
	// Channel flag waiting for signal to close connection.
	sentinel chan bool
	// Routing hash table eg. {Socket: Conn interface}.
	router *router
	// Pubsub notifications.
	events *events
	// Configuration settings
	config Config
}

// New create a new node with default
func New(config Config) *Node {
	return &Node{
		make(chan bool),
		newRouter(),
		newEvents(),
		config,
	}
}

// Signals proxy channels to subscriber.
// The listening routine should be stopped using context param.
func (n *Node) Signals(ctx context.Context) <-chan Signal {
	ch := make(chan Signal)
	go n.events.Listen(ctx, ch)
	return ch // read only channel for raw messages
}

// Send emit a new message using peer id.
// If peer id doesn't exists or peer is not connected return error.
// Calling Send extends write deadline.
func (n *Node) Send(id ID, message []byte) (int, error) {
	// Check if id exists in connected peers
	peer := n.router.Query(id)
	if peer == nil {
		err := fmt.Errorf("remote peer disconnected: %s", id.String())
		return 0, errSendingMessage(err)
	}

	bytes, err := peer.Send(message)
	// An idle timeout can be implemented by repeatedly extending
	// the deadline after successful Read or Write calls.
	// SetWriteDeadline sets the deadline for future Write calls
	// and any currently-blocked Write call.
	// Even if write times out, it may return n > 0, indicating that
	// some of the data was successfully written.
	idle := futureDeadLine(n.config.PeerDeadline())
	peer.SetWriteDeadline(idle)
	return bytes, err
}

// watch keep running waiting for incoming messages.
// After every new message the connection is verified, if local connection is closed or remote peer is disconnected the routine is stopped.
// Incoming message monitor is suggested to be processed in go routines.
func (n *Node) watch(peer *peer) {

KEEPALIVE:
	for {

		// Waiting for new incoming message
		buf, err := peer.Listen(n.config.MaxPayloadSize())
		// If connection is closed
		if n.Closed() {
			// stop routines watching for peers
			return
		}

		// OverflowError is returned when the incoming payload exceed the expected size
		_, overflow := err.(OverflowError)

		// Don't stop listening for peer if overflow payload is returned.
		if err != nil && !overflow {
			// net: don't return io.EOF from zero byte reads
			// Notify about the remote peer state
			n.events.PeerDisconnected(peer)
			// Remove peer from router table
			n.router.Remove(peer)
			return
		}

		if buf == nil {
			log.Printf("buffer nil with err: %v", err)
			// `buf` is nil if no more bytes received but peer is still connected
			// Keep alive always that zero bytes are not received
			break KEEPALIVE
		}

		// Emit new incoming message notification
		n.events.NewMessage(peer, buf)
		// An idle timeout can be implemented by repeatedly extending
		// the deadline after successful Read or Write calls.
		// SetReadDeadline sets the deadline for future Read calls
		// and any currently-blocked Read call.
		idle := futureDeadLine(n.config.PeerDeadline())
		peer.SetReadDeadline(idle)

	}

}

// routing initialize route in routing table from connection interface.
// If TCP protocol is used connection is enforced to keep alive.
// Return err if max peers connected exceed MaxPeerConnected otherwise return new peer added to table.
func (n *Node) routing(conn net.Conn) (*peer, error) {

	// Assertion for tcp connection to keep alive
	connection, isTCP := conn.(*net.TCPConn)
	if isTCP {
		// If tcp enforce keep alive connection
		// SetKeepAlive sets whether the operating system should send keep-alive messages on the connection.
		connection.SetKeepAlive(true)
	}

	// Drop connections if max peers exceeded
	if n.router.Len() >= n.config.MaxPeersConnected() {
		log.Printf("max peers exceeded: MaxPeerConnected = %d", n.config.MaxPeersConnected())
		return nil, errExceededMaxPeers(n.config.MaxPeersConnected())
	}

	// Initial deadline for connection.
	// A deadline is an absolute time after which I/O operations
	// fail instead of blocking. The deadline applies to all future
	// and pending I/O, not just the immediately following call to
	// Read or Write. After a deadline has been exceeded, the
	// connection can be refreshed by setting a deadline in the future.
	// ref: https://pkg.go.dev/net#Conn
	idle := futureDeadLine(n.config.PeerDeadline())
	connection.SetDeadline(idle)
	// We need to know how interact with peer based on socket and connection
	peer := newPeer(connection)
	// Store new peer in router table
	n.router.Add(peer)
	return peer, nil
}

// Listen start listening on the given address and wait for new connection.
// Return error if error occurred while listening.
func (n *Node) Listen() error {

	addr := n.config.SelfListeningAddress()
	protocol := n.config.Protocol()

	listener, err := net.Listen(protocol, addr)
	if err != nil {
		return err
	}

	log.Printf("listening on %s", addr)
	//wait until sentinel channel is closed to close listener
	defer func() {
		err := listener.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
	}()

	for {
		// Block/Hold while waiting for new incoming connection
		// Synchronized incoming connections
		conn, err := listener.Accept()
		// If connection is closed
		if n.Closed() {
			// Graceful stop listening
			return nil
		}

		if err != nil {
			log.Printf("error accepting connection: %v", err)
			return errBindingConnection(err)
		}

		// Routing for accepted connection
		// TODO 1 - accepted?
		// TODO 2 - Run handshake?
		// TODO 3 - OK? then add peer to router newPeer(id, securedConn)
		// TODO 4 - then trigger event and start watching for peer
		peer, err := n.routing(conn)
		if err != nil {
			conn.Close() // Drop connection
			continue
		}
		// Wait for incoming messages
		// This routine will stop when Close() is called
		go n.watch(peer)
		// Dispatch event for new peer connected
		n.events.PeerConnected(peer)
	}

}

// Closed check connection state.
// Return true for connection open else false.
func (n *Node) Closed() bool {
	select {
	// select await for sentinel if not closed then default is returned.
	case <-n.sentinel:
		return true
	default:
		return false
	}
}

// Close all peers connections and stop listening.
func (n *Node) Close() {
	for _, p := range n.router.Table() {
		go func(peer *peer) {
			if err := peer.Close(); err != nil {
				log.Printf("error when shutting down connection: %v", err)
			}
		}(p)
	}

	// flush connected peers
	n.router.Flush()
	// If channel get closed then all routines waiting for connections
	// or waiting for incoming messages get closed too.
	close(n.sentinel)
}

// Dial attempt to connect to remote node and add connected peer to routing table.
// Return error if error occurred while dialing node.
func (n *Node) Dial(addr string) error {

	protocol := n.config.Protocol()   // eg. tcp
	timeout := n.config.DialTimeout() // max time waiting for dial.

	// Start dialing to address
	conn, err := net.DialTimeout(protocol, addr, timeout)
	log.Printf("dialing to %s", addr)

	if err != nil {
		return errDialingNode(err)
	}

	// Routing for dialed connection
	peer, err := n.routing(conn)
	if err != nil {
		conn.Close() // Drop connection
		return errDialingNode(err)
	}

	// Wait for incoming messages
	// This routine will stop when Close() is called
	go n.watch(peer)
	// Dispatch event for new peer connected
	n.events.PeerConnected(peer)
	return nil
}
