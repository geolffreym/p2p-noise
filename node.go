//Copyright (c) 2022, Geolffrey Mena <gmjun2000@gmail.com>

//P2P Noise Secure handshake.
//
//See also: http://www.noiseprotocol.org/noise.html#introduction
package noise

import (
	"context"
	"io"
	"log"
	"net"

	"github.com/geolffreym/p2p-noise/conf"
	"github.com/geolffreym/p2p-noise/errors"
)

// Default protocol
const PROTOCOL = "tcp"

type Node struct {
	sentinel chan bool      // Channel flag waiting for signal to close connection.
	router   *Router        // Routing hash table eg. {Socket: Conn interface}.
	events   *Events        // Pubsub notifications.
	settings *conf.Settings // Configuration settings
}

func NewNode(c ...conf.Setting) *Node {
	// Create settings from params and write in settings reference
	settings := conf.NewSettings()
	settings.Write(c...)

	return &Node{
		router:   newRouter(),
		events:   newEvents(),
		sentinel: make(chan bool),
		settings: settings,
	}
}

// Events proxy channels to subscriber.
// The listening routine should be stopped using context param.
func (n *Node) Events(ctx context.Context) <-chan Message {
	ch := make(chan Message)
	go n.events.Subscriber().Listen(ctx, ch)
	return ch // read only channel <-chan
}

// watch keep running waiting for incoming messages.
// After every new message the connection is verified, if local connection is closed or remote peer is disconnected the watch routine is stopped.
// Incoming message monitor is suggested to be processed in go routines.
func (n *Node) watch(peer *Peer) {
	// Recycle memory buffer
	buf := make([]byte, 1024)

KEEPALIVE:
	for {
		// Sync buffered IO reading
		_, err := peer.Read(buf)
		// If connection is closed
		if n.Closed() {
			// stop routines watching for peers
			return
		}

		if err != nil {
			// net: don't return io.EOF from zero byte reads
			// if err == io.EOF then peer connection is closed
			_, isNetError := err.(*net.OpError)
			if err == io.EOF || isNetError {
				// Close disconnected peer
				if err := peer.Close(); err != nil {
					log.Fatal(errors.Closing(err).Error())
				}

				// Notify about the remote peer state
				n.events.PeerDisconnected([]byte(peer.Socket()))
				// Remove peer from router table
				n.router.Remove(peer)
				return
			}

			// Keep alive always that zero bytes are not received
			break KEEPALIVE
		}

		// Emit new incoming message notification
		n.events.NewMessage(buf)
	}

}

// routing initialize route in routing table from connection interface.
// If TCP protocol is used connection is enforced to keep alive.
// It return new peer added to table.
func (n *Node) routing(conn net.Conn) (*Peer, error) {

	// Assertion for tcp connection to keep alive
	connection, isTCP := conn.(*net.TCPConn)
	if isTCP {
		// If tcp enforce keep alive connection
		// SetKeepAlive sets whether the operating system should send keep-alive messages on the connection.
		connection.SetKeepAlive(true)
	}

	// Drop connections if max peers exceeded
	if n.router.Len() >= n.settings.MaxPeersConnected {
		log.Fatalf("max peers exceeded: MaxPeerConnected = %d", n.settings.MaxPeersConnected)
		return nil, errors.Exceeded(n.settings.MaxPeersConnected)
	}

	// Routing connections
	remote := connection.RemoteAddr().String()
	// eg. 192.168.1.1:8080
	socket := Socket(remote)
	// We need to know how interact with peer based on socket and connection
	peer := newPeer(socket, connection)
	n.router.Add(peer)
	return peer, nil
}

// Listen start listening on the given address and wait for new connection.
// Return error if error occurred while listening.
func (n *Node) Listen(addr string) error {

	listener, err := net.Listen(PROTOCOL, addr)
	if err != nil {
		return err
	}

	// Dispatch event on start listening
	n.events.Listening([]byte(addr))
	//wait until sentinel channel is closed to close listener
	go func(listener net.Listener) {
		<-n.sentinel
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
		if n.Closed() {
			// Graceful stop listening
			return nil
		}

		if err != nil {
			log.Fatal(errors.Binding(err).Error())
			return err
		}

		// Routing for accepted connection
		peer, err := n.routing(conn)
		if err != nil {
			conn.Close() // Drop connection
			continue
		}

		go n.watch(peer) // Wait for incoming messages
		// Dispatch event for new peer connected
		payload := []byte(peer.Socket())
		n.events.PeerConnected(payload)
	}

}

// Table return current routing table.
func (n *Node) Table() Table {
	return n.router.Table()
}

// Closed check connection state.
// Return true for connection open else false.
func (n *Node) Closed() bool {
	select {
	case <-n.sentinel:
		return true
	default:
		return false
	}
}

// Close all peers connections and stop listening
func (n *Node) Close() {
	for _, peer := range n.router.Table() {
		go func(p *Peer) {
			if err := p.Close(); err != nil {
				log.Fatal(errors.Closing(err).Error())
			}
		}(peer)
	}

	// Dispatch event on node get closed
	n.events.ClosedConnection()
	// If channel get closed then all routines waiting for connections
	// or waiting for incoming messages get closed too.
	close(n.sentinel)
}

// Dial attempt to connect to remote node and add connected peer to routing table.
// Return error if error occurred while dialing node.
func (n *Node) Dial(addr string) error {
	conn, err := net.Dial(PROTOCOL, addr)
	if err != nil {
		return errors.Dialing(err, addr)
	}

	// Routing for dialed connection
	peer, err := n.routing(conn)
	if err != nil {
		conn.Close() // Drop connection
		return errors.Dialing(err, addr)
	}

	go n.watch(peer) // Wait for incoming messages
	// Dispatch event for new peer connected
	n.events.PeerConnected([]byte(peer.Socket()))
	return nil
}
