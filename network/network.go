package network

import (
	"fmt"
	"io"
	"log"
	"net"
)

const PROTOCOL = "tcp"

type Network struct {
	events Events
	table  Router
	closed chan bool
}

func New() *Network {
	return &Network{
		table:  make(Router),
		events: make(Events),
		closed: make(chan bool, 1),
	}
}

// Create a route from net connection
func (network *Network) route(conn net.Conn) *Route {
	return &Route{
		conn:   conn,
		socket: Socket(conn.RemoteAddr().String()),
	}
}

// Initialize route in routing table
func (network *Network) routing(conn net.Conn) *Route {
	// Keep routing for each connection
	socket := Socket(conn.RemoteAddr().String())
	route := network.route(conn)
	network.table.Add(socket, route)
	return route
}

// Run routed stream message in goroutine
// Each incoming message processed in concurrent approach
func (network *Network) stream(route *Route) {
	go func(n *Network, r *Route) {
		buf := make([]byte, 1024)

	KEEPALIVE:
		for {
			// Stop routine
			if n.IsClosed() {
				return
			}

			_, err := r.Read(buf)
			if err != nil {
				if err == io.EOF {
					break KEEPALIVE
				}
			}

			// TODO Need refactor to handle biggest messages
			// Emit new incoming
			n.events.Emit(MESSAGE, r, buf)

		}
	}(network, route)
}

// Concurrent Bind network and set routing to start listening for streams
func (network *Network) bind(listener net.Listener) {
	go func(n *Network, l net.Listener) {
		for {
			// Block/Hold while waiting for new incoming connection
			conn, err := l.Accept()
			if err != nil || n.IsClosed() {
				log.Fatalf("connection closed or cannot be established: %v", err)
				return
			}

			// Routing for connection
			route := n.routing(conn)
			n.stream(route)
			// Dispatch event
			n.events.Emit(NEWPEER, route)
		}
	}(network, listener)
}

// Start listening on the given address and wait for new connection
func (network *Network) Listen(addr string) (*Network, error) {
	listener, err := net.Listen(PROTOCOL, addr)
	if err != nil {
		return nil, fmt.Errorf("error trying to listen on %s: %v", addr, err)
	}

	// Concurrent processing for each incoming connection
	network.bind(listener)
	// Dispatch event on start listening
	network.events.Emit(LISTENING, &Route{})
	return network, nil
}

func (network *Network) Table() Router {
	return network.table
}

// Non-blocking check connection state
func (network *Network) IsClosed() bool {
	select {
	case <-network.closed:
		return true
	default:
		return false
	}
}

// Close all peers connections
func (network *Network) Close() {
	for _, route := range network.table {
		go func(r *Route) {
			if err := r.Close(); err != nil {
				log.Fatalf("error when shutting down connection: %s", err)
			}
		}(route)
	}

	network.closed <- true
}

// Dial to a network node and add route to table
func (network *Network) Dial(addr string) (*Network, error) {
	conn, err := net.Dial(PROTOCOL, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %v", addr, err)
	}

	// Routing for connection
	route := network.routing(conn)
	network.stream(route)
	// Dispatch event
	network.events.Emit(NEWPEER, route)
	return network, nil
}

// Set handler to handle incoming messages
func (network *Network) AddEventListener(event Event, h Handler) *Network {
	network.events.AddListener(event, h)
	return network
}
