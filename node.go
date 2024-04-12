// Package noise implements the Noise Protocol for peer-to-peer communication.
// For more information about the Noise Protocol, please visit: [Noise Protocol].
//
// [Noise Protocol]: http://www.noiseprotocol.org/noise.html
package noise

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/oxtoacart/bpool"
)

// futureDeadline calculate and return a new time for deadline since now.
func futureDeadLine(deadline time.Duration) time.Time {
	if deadline == 0 {
		// deadline 0 = no deadline
		return time.Time{}
	}

	// how long should i wait for activity?
	// since now add a new future deadline
	return time.Now().Add(deadline * time.Second)
}

type Config interface {
	// Default "tcp"
	Protocol() string
	// Default 0.0.0.0:8010
	SelfListeningAddress() string
	// Default 0
	Linger() int
	// Default 100
	MaxPeersConnected() uint8
	// Default 10 << 20 = 10MB
	PoolBufferSize() int
	// Default 0
	IdleTimeout() time.Duration
	// Default 5 seconds
	DialTimeout() time.Duration
	// Default 1800 seconds
	KeepAlive() time.Duration
}

// Node represents a network node capable of handling connections,
// routing messages, and managing configurations.
type Node struct {
	// Bound local network listener.
	listener net.Listener
	// Routing hash table eg. {Socket: Conn interface}.
	router *router
	// Pubsub notifications.
	events *events
	// Global buffer pool
	pool BytePool
	// Configuration settings
	config Config
}

// New create a new node with defaults
func New(config Config) *Node {
	// Max allowed "pools" is related to max active peers.
	maxPools := int(config.MaxPeersConnected())
	// Width of global pool buffer
	maxBufferSize := config.PoolBufferSize()
	pool := bpool.NewBytePool(maxPools, maxBufferSize)

	return &Node{
		router: newRouter(),
		events: newEvents(),
		pool:   pool,
		config: config,
	}
}

// Signals initiates the signaling process to proxy channels to subscribers.
// It returns a channel of type Signal to intercept events and a cancel function to stop the listening routine.
// The channel is closed during the cancellation of listening.
func (n *Node) Signals() (<-chan Signal, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	// this channel is closed during listening cancellation
	ch := make(chan Signal)
	// forward signals to internal events signaling
	go n.events.Listen(ctx, ch)
	return ch, cancel // read only channel for raw messages
}

// Disconnect close all the peer connections without stop listening.
func (n *Node) Disconnect() {
	log.Print("closing connections and shutting down node..")
	for peer := range n.router.Table() {
		if err := peer.Close(); err != nil {
			log.Printf("error when shutting down connection: %v", err)
		}
	}
}

// Send emits a new message using a peer ID.
// It returns the total bytes sent if there is no error; otherwise, it returns 0.
// If the peer ID doesn't exist or the peer is not connected, it returns an error.
// Calling Send extends the write deadline.
func (n *Node) Send(rawID string, message []byte) (uint32, error) {
	id := newIDFromString(rawID)
	// Check if id exists in connected peers
	// check in-band error
	peer, ok := n.router.Query(id)
	if !ok {
		err := fmt.Errorf("remote peer disconnected: %s", id.String())
		return 0, errSendingMessage(err)
	}

	bytes, err := peer.Send(message)
	// An idle timeout can be implemented by repeatedly extending
	// the deadline after successful Read or Write calls.
	idle := futureDeadLine(n.config.IdleTimeout())
	peer.SetDeadline(idle)
	return bytes, err
}

// watch keeps running, waiting for incoming messages.
// After receiving each new message, the connection is verified. If the local connection is closed or the remote peer is disconnected, the routine stops.
// It is suggested to process incoming messages in separate goroutines.
func (n *Node) watch(peer *peer) {

KEEPALIVE:
	for {

		// Waiting for new incoming message
		buf, err := peer.Listen()
		if err != nil {
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

		log.Print("receiving message from remote")
		// Emit new incoming message notification
		n.events.NewMessage(peer, buf)
		// An idle timeout can be implemented by repeatedly extending
		// the deadline after successful Read or Write calls.
		idle := futureDeadLine(n.config.IdleTimeout())
		peer.SetDeadline(idle)

	}

}

// setupTCPConnection configures the behavior of a TCP connection.
// It takes a net.TCPConn connection and modifies its settings according to the Node configuration.
// If any of the configurations cannot be fulfilled, it returns an error.
func (n *Node) setupTCPConnection(conn *net.TCPConn) error {
	// If tcp enforce keep alive connection.
	// SetKeepAlive sets whether the operating system should send keep-alive messages on the connection.
	// We can modify the behavior of the connection using idle timeout and keep alive or just disable keep alive and force the use of idle timeout to determine the inactivity of remote nodes.
	// ref: https://support.f5.com/csp/article/K13004262
	if err := conn.SetKeepAlivePeriod(n.config.KeepAlive()); err != nil {
		return err
	}
	// Set linger time in seconds to wait to discard unsent data after close.
	// discard after N seconds unsent packages on close connection.
	if err := conn.SetLinger(n.config.Linger()); err != nil {
		return err
	}

	return nil
}

// handshake initiates a new handshake for an incoming or dialed connection.
// After the handshake completes, a new session is created, and a new peer is added to the router.
// If the TCP protocol is used, the connection is enforced to keep alive.
// Returns an error if the maximum number of connected peers exceeds MaxPeersConnected; otherwise, returns nil.
func (n *Node) handshake(conn net.Conn, initialize bool) error {

	// Assertion for tcp connection to keep alive
	log.Print("starting handshake")
	connection, isTCP := conn.(*net.TCPConn)
	if isTCP {
		// Setup network parameters to control connection behavior.
		if err := n.setupTCPConnection(connection); err != nil {
			return errSettingUpConnection(err)
		}
	}

	// Drop connections if max peers exceeded
	if n.router.Len() >= n.config.MaxPeersConnected() {
		connection.Close() // Drop connection :(
		log.Printf("max peers exceeded: MaxPeerConnected = %d", n.config.MaxPeersConnected())
		return errExceededMaxPeers(n.config.MaxPeersConnected())
	}

	// Stage 1 -> run handshake
	h, err := newHandshake(connection, initialize)
	if err != nil {
		log.Printf("error while creating handshake: %s", err)
		return err
	}

	err = h.Start() // start the handshake
	if err != nil {
		log.Printf("error while starting handshake: %s", err)
		return err
	}

	// Stage 2 -> get a secure session
	// All good with handshake? Then get a secure session.
	log.Print("handshake complete")
	session := h.Session()
	// Stage 3 -> create a peer and add it to router
	// Routing for secure session
	peer := n.routing(session)
	// Keep watching for incoming messages
	// This routine will stop when Close() is called
	go n.watch(peer)
	// Dispatch event for new peer connected
	n.events.PeerConnected(peer)
	return nil
}

// routing initializes a route in the routing table from a session.
// It returns the recently added peer.
func (n *Node) routing(conn *session) *peer {
	// Initial deadline for connection.
	// A deadline is an absolute time after which I/O operations
	// fail instead of blocking. The deadline applies to all future
	// and pending I/O, not just the immediately following call to
	// Read or Write. After a deadline has been exceeded, the
	// connection can be refreshed by setting a deadline in the future.
	// ref: https://pkg.go.dev/net#Conn
	idle := futureDeadLine(n.config.IdleTimeout())
	conn.SetDeadline(idle)
	// We need to know how interact with peer based on socket and connection
	peer := newPeer(conn)
	// Bind global buffer pool to peer.
	// Pool buffering reduce memory allocation latency.
	peer.BindPool(n.pool)
	// Store new peer in router table
	n.router.Add(peer)
	return peer
}

// LocalAddr returns the local address assigned to node.
// If node is not listening nil is returned instead.
func (n *Node) LocalAddr() net.Addr {
	if n.listener == nil {
		return nil
	}

	return n.listener.Addr()
}

// Listen start listening on the given address and wait for new connection.
// Return error if error occurred while listening.
func (n *Node) Listen() error {

	addr := n.config.SelfListeningAddress() // eg. 0.0.0.0
	protocol := n.config.Protocol()         // eg. tcp
	listener, err := net.Listen(protocol, addr)
	log.Printf("listening on %s", addr)

	if err != nil {
		log.Printf("error listening on %s: %v", addr, err)
		return err
	}

	// The order here is IMPORTANT.
	// We set listener first then we notify listening event, otherwise a race condition is caused.
	n.listener = listener        // keep reference to current listener.
	n.events.SelfListening(addr) // emit listening event

	for {
		// Block/Hold while waiting for new incoming connection
		// Synchronized incoming connections
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("error accepting connection %s", err)
			return errBindingConnection(err)
		}

		// Run handshake for incoming connection
		// We need to run in a separate goroutine to improve time performance between nodes requesting connections.
		go n.handshake(conn, false)
	}

}

// Close all peers connections and stop listening.
func (n *Node) Close() error {

	// close peer connections
	go n.Disconnect()

	// stop listener for listening node only
	if n.listener == nil {
		return nil
	}

	if err := n.listener.Close(); err != nil {
		return err
	}

	return nil
}

// Dial attempts to connect to a remote node and adds the connected peer to the routing table.
// It returns an error if an error occurred while dialing the node.
func (n *Node) Dial(addr string) error {
	protocol := n.config.Protocol()   // eg. tcp
	timeout := n.config.DialTimeout() // max time waiting for dial.

	// Start dialing to address
	conn, err := net.DialTimeout(protocol, addr, timeout)
	log.Printf("dialing to %s", addr)

	if err != nil {
		return errDialingNode(err)
	}

	// Run handshake for dialed connection
	if err = n.handshake(conn, true); err != nil {
		return err
	}

	return nil
}
