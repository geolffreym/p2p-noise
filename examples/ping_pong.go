//Copyright (c) 2022, Geolffrey Mena <gmjun2000@gmail.com>

//P2P Noise Secure handshake.
//
//See also: http://www.noiseprotocol.org/noise.html#introduction
package main

import (
	"context"
	"log"
	"os"

	noise "github.com/geolffreym/p2p-noise"
	"github.com/geolffreym/p2p-noise/config"
)

func main() {

	// Create configuration from params and write in configuration reference
	configuration := config.New()
	configuration.Write(
		config.SetMaxPeersConnected(10),
		config.SetPeerDeadline(1800),
	)

	// Node factory
	args := os.Args[1:]
	ip := args[0]
	port := args[1]

	remote := noise.Socket(ip + ":" + port)
	node := noise.New(configuration)

	// Network events channel
	ctx, cancel := context.WithCancel(context.Background())
	var signals <-chan noise.SignalCtx = node.Signals(ctx)

	go func() {
		// Wait for incoming message channel.
		for signal := range signals {
			// Here could be handled events
			switch signal.Type() {
			case noise.NewPeerDetected:
				// When a new peer is connected. Start ping pong game.
				log.Printf("New Peer connected: %s \n", signal.Payload())
				signal.Reply([]byte("ping")) // start game

			case noise.MessageReceived:
				// When we receive a message, check the content message and reply "ping" or "pong"
				message := string(signal.Payload())
				log.Printf("New Message: %s", message)
				if message == "ping" {
					signal.Reply([]byte("pong"))
				} else {
					signal.Reply([]byte("ping"))
				}

			case noise.PeerDisconnected:
				// What we do when a peer get disconnected?
				log.Printf("Peer disconnected")
				cancel() // stop listening for events
			}
		}
	}()

	// ... some code here
	err := node.Dial(remote)
	if err != nil {
		log.Fatal(err)
	}
	// node.Close()

	// ... more code here
	node.Listen()

}
