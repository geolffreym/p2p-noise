package network

import (
	"fmt"
	"log"
	"net"
)

const PROTOCOL = "tcp"

type Network struct {
	Addr   string
	events Events
	table  Router
}

func New(addr string) *Network {
	return &Network{
		Addr:   addr,
		table:  make(Router),
		events: make(Events),
	}
}

// Start listening on the given address
func (network *Network) Listen() (*Network, error) {
	listener, err := net.Listen(PROTOCOL, network.Addr)
	if err != nil {
		return nil, fmt.Errorf("error trying to listen on %s: %v", network.Addr, err)
	}

	// Concurrent processing for each incoming connection
	go func(n *Network, l net.Listener) {
		for {
			// Block/Hold while waiting for new incoming connection
			conn, err := l.Accept()
			if err != nil {
				log.Fatalf("connection closed or cannot be established: %v", err)
				return
			}

			// Routing for connection
			route := n.routing(conn)
			n.stream(route)
			// Dispatch event
			network.events.Emit(NEWPEER, route)
		}
	}(network, listener)

	// Dispatch event
	network.events.Emit(LISTENING, &Route{})
	return network, nil
}

func (network *Network) Table() Router {
	return network.table
}

// Close all peers connections
func (network *Network) Close() {
	for _, route := range network.table {
		route.Close()
	}
}

// Dial to a network node
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
func (network *Network) stream(route *Route) {
	// Each incoming message processed in concurrent approach
	go func(n *Network, r *Route) {
		buf := make([]byte, 1024)
		for {

			_, err := r.Read(buf)
			if err != nil {
				continue
			}

			// Emit new incoming
			n.events.Emit(MESSAGE, r, buf)
		}

	}(network, route)
}

// Set handler to handle incoming messages
func (network *Network) AddEventListener(event Event, h Handler) *Network {
	network.events.AddListener(event, h)
	return network
}
