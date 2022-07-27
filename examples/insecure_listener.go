//Copyright (c) 2022, Geolffrey Mena <gmjun2000@gmail.com>

//P2P Noise Secure handshake.
//
//See also: http://www.noiseprotocol.org/noise.html#introduction
package main

import (
	"context"
	"log"

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
	node := noise.New(configuration)
	// Network events channel
	ctx, cancel := context.WithCancel(context.Background())
	var signals <-chan noise.SignalCtx = node.Signals(ctx)

	go func() {
		for signal := range signals {
			// Here could be handled events
			if signal.Type() == noise.NewPeerDetected {
				log.Printf("New Peer connected: %s \n", signal.Payload())
				signal.Reply([]byte("ping")) // start game
			}

			if signal.Type() == noise.MessageReceived {
				message := string(signal.Payload())

				if message == "ping" {
					signal.Reply([]byte("pong"))
				} else {
					signal.Reply([]byte("ping"))
				}
			}

			if signal.Type() == noise.PeerDisconnected {
				log.Printf("Peer disconnected")
				cancel() // stop listening for events
			}
		}
	}()

	// ... some code here
	err := node.Dial("192.168.1.17:8010")
	if err != nil {
		log.Fatal(err)
	}
	// node.Close()

	// ... more code here
	node.Listen()

}
