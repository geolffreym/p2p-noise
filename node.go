// Network implements a lightweight TCP communication.
// Offers pretty basic features to communicate between nodes.
//
// See also: https://pkg.go.dev/net#Conn
package noise

import (
	"context"
	"io"
	"log"
	"net"

	"github.com/geolffreym/p2p-noise/errors"
)

// Default protocol
const PROTOCOL = "tcp"

type (
	Socket = string
	Table  = map[Socket]*Peer
)

// Node implement communication logic
type Node struct {
	sentinel chan bool // Channel flag waiting for signal to close connection.
	router   *Router   // Routing hash table eg. {Socket: Conn interface}.
	events   *Events   // Pubsub notifications.
}

// Node factory
// It receive a param events message handler for network.
func NewNode() *Node {
	return &Node{
		router:   newRouter(),
		events:   newMessenger(),
		sentinel: make(chan bool),
	}
}

// Accessor to messenger events listener
func (node *Node) Intercept(cb Observer) context.CancelFunc {
	return node.events.Listen(cb)
}

// watch watchdog for incoming messages.
// incoming message monitor is suggested to be processed in go routines.
func (node *Node) watch(peer *Peer) {
	buf := make([]byte, 1024)

KEEPALIVE:
	for {
		// Sync buffer reading
		_, err := peer.Receive(buf)
		// If connection is closed
		// stop routines watching peers
		if node.Closed() {
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

				//Notify to node about the peer state
				node.events.PeerDisconnected([]byte(peer.Socket()))
				// Remove peer from router table
				node.router.Delete(peer)
				return
			}

			// Keep alive always that zero bytes are not received
			break KEEPALIVE
		}

		// Emit new incoming message notification
		node.events.NewMessage(buf)
	}

}

// routing initialize route in routing table from connection interface
// Return new peer added to table
func (node *Node) routing(conn net.Conn) *Peer {

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
	peer := newPeer(socket, connection)
	node.router.Add(peer)
	return peer
}

// Listen start listening on the given address and wait for new connection.
// Return error if error occurred while listening.
func (node *Node) Listen(addr string) error {
	listener, err := net.Listen(PROTOCOL, addr)
	if err != nil {
		return err
	}

	// Dispatch event on start listening
	node.events.Listening([]byte(addr))
	// monitor connection to close listener
	go func(listener net.Listener) {
		<-node.sentinel
		err := listener.Close()
		if err != nil {
			log.Fatal(errors.Closing(err).Error())
		}
	}(listener)

	for {
		// Block/Hold while waiting for new incoming connection
		// Synchronized incoming connections
		conn, err := listener.Accept()
		// If connection is closed
		// Graceful stop listening
		if node.Closed() {
			return nil
		}

		if err != nil {
			log.Fatal(errors.Binding(err).Error())
			return err
		}

		peer := node.routing(conn) // Routing for connection
		go node.watch(peer)        // Wait for incoming messages
		// Dispatch event for new peer connected
		payload := []byte(peer.Socket())
		node.events.PeerConnected(payload)
	}

}

// Return current routing table
func (node *Node) Table() Table {
	return node.router.Table()
}

// Closed Non-blocking check connection state.
// Return true for connection open else false
func (node *Node) Closed() bool {
	select {
	case <-node.sentinel:
		return true
	default:
		return false
	}
}

// Close all peers connections and stop listening
func (node *Node) Close() {
	for _, peer := range node.router.Table() {
		go func(p *Peer) {
			if err := p.Close(); err != nil {
				log.Fatal(errors.Closing(err).Error())
			}
		}(peer)
	}

	// Dispatch event on node get closed
	node.events.ClosedConnection()
	// If channel get closed then all routines waiting for connections
	// or waiting for incoming messages get closed too.
	close(node.sentinel)
}

// Dial to node and add connected peer to routing table
// Return error if error occurred while dialing node.
func (node *Node) Dial(addr string) error {
	conn, err := net.Dial(PROTOCOL, addr)
	if err != nil {
		return errors.Dialing(err, addr)
	}

	peer := node.routing(conn) // Routing for connection
	go node.watch(peer)        // Wait for incoming messages
	// Dispatch event for new peer connected
	node.events.PeerConnected([]byte(peer.Socket()))
	return nil
}
