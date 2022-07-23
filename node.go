//Copyright (c) 2022, Geolffrey Mena <gmjun2000@gmail.com>
//
//P2P Noise Secure handshake.
//See also: http://www.noiseprotocol.org/noise.html#introduction
package noise

import (
	"context"
	"time"
)

type Config interface {
	MaxPeersConnected() uint8
	PeerDeadline() time.Duration
}

type Node struct {
	//monitor
	net *network // networking
	// sec *Noise   // security

}

func New(config Config) *Node {
	return &Node{
		newNetwork(config),
	}
}

// Events its a middleware/proxy to intercept and pre-process raw incoming messages
func (n *Node) Events(ctx context.Context) <-chan Message {
	events := n.net.Events(ctx)
	ch := make(chan Message)

	// Intercept messages to split, public from internal events
	go func() {
		for msg := range events {
			// Here could be handled events
			if msg.Type() == NewPeerDetected {
				// TODO start handshake?
				// cancel() // stop listening for events

			}

			if msg.Type() == MessageReceived {
				// Decrypt message
			}

			// Redirect message
			ch <- msg
		}
	}()

	return ch
}

func (n *Node) SendMessage() {
	// Encrypt message
}

func (n *Node) Handshake() {

}

// Listen forward method using default node address
func (n *Node) Listen(addr Socket) error {
	return n.net.Listen(addr)
}
