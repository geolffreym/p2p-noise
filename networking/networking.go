package networking

import (
	"fmt"
	"io"
	"log"
	"net"
)

const PROTOCOL = "tcp"
const BUFFER_SIZE = 1024

type Network struct {
	Addr     string
	Port     uint16
	listener net.Listener
}

type Connection struct {
	conn net.Conn
}

func New(addr string, port uint16) *Network {
	return &Network{
		Addr: addr,
		Port: port,
	}
}

// Start listening on the given address
func (network *Network) Listen() *Network {

	var err error
	nodeAddress := fmt.Sprintf("%s:%d", network.Addr, network.Port)
	network.listener, err = net.Listen(PROTOCOL, nodeAddress)

	if err != nil {
		panic("Cannot start listening on " + nodeAddress)
	}

	return network

}

// Start "handshake" for incoming packages
func (network *Network) RunHandshake() *Network {

	// Concurrent processing for each incoming message
	go func(listener net.Listener) {
		for {
			// Block/Hold while waiting for new incoming connection + handshake
			connection, err := listener.Accept()
			if err != nil {
				log.Fatalf("Connection closed or cannot be established: %v", err)
				return
			}

			// Each incoming message processed in concurrent approach
			go network.HandleMessage(connection)

		}
	}(network.listener)

	return network
}

// Each incoming message needs to be handled
// Possible interaction with node:
// Register | Message
func (n *Network) HandleMessage(conn net.Conn) {
	defer conn.Close()
	message := make([]byte, BUFFER_SIZE)

	for {
		b, err := conn.Read(message)
		// If not more bytes to read or error
		if b == 0 || err != nil {
			// If error is not related to end of message
			if err != io.EOF {
				log.Fatalf("Error reading message: %v", err)
			}
		}

		log.Print(message)
	}
}
