package network

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

const PROTOCOL = "tcp"

type Network struct {
	Addr     string
	handler  handler
	table    router
	listener net.Listener
}

func New(addr string) *Network {
	return &Network{
		Addr:  addr,
		table: make(router),
	}
}

// Start listening on the given address
func (network *Network) Listen() (*Network, error) {
	listener, err := net.Listen(PROTOCOL, network.Addr)
	if err != nil {
		return nil, fmt.Errorf("error trying to listen on %s: %v", network.Addr, err)
	}

	// Node setup
	log.Printf("listening on %s", network.Addr)
	network.listener = listener
	return network, nil
}

// Close all peers connections
func (network *Network) Close() {
	for socket, route := range network.table {
		log.Printf("Closing connection for peer %s", socket)
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
	log.Printf("Connecting to network %s", addr)
	route := network.routing(conn)
	network.stream(route)
	return network, nil
}

// Create a route from net connection
func (network *Network) route(conn net.Conn) *Route {
	return &Route{
		conn:   conn,
		socket: socket(conn.RemoteAddr().String()),
		stream: bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)),
	}
}

// Initialize route in routing table
func (network *Network) routing(conn net.Conn) *Route {
	// Keep routing for each connection
	socket := socket(conn.RemoteAddr().String())
	route := network.route(conn)
	network.table.Add(socket, route)
	return route
}

// Run routed stream connection in goroutine concurrency
func (network *Network) stream(route *Route) {
	if network.handler == nil {
		log.Fatalf("invalid handler for streaming connection")
		return
	}

	// Each incoming message processed in concurrent approach
	go func(n *Network, r *Route) {
		n.handler(route)
	}(network, route)
}

// Set handler to handle incoming messages
func (network *Network) SetHandler(h handler) *Network {
	network.handler = h
	return network
}

// Bind network in thread to handle in/out streams
func (network *Network) Bind() *Network {
	// Concurrent processing for each incoming connection
	go func(listener net.Listener) {
		for {
			// Block/Hold while waiting for new incoming connection
			conn, err := listener.Accept()
			if err != nil {
				log.Fatalf("connection closed or cannot be established: %v", err)
				return
			}

			// Routing for connection
			log.Printf("accepted connection from %v", conn.RemoteAddr())
			route := network.routing(conn)
			network.stream(route)
		}
	}(network.listener)

	return network
}
