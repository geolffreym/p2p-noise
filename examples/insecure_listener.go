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
	var events <-chan noise.Message = node.Events(ctx)

	go func() {
		for msg := range events {
			// Here could be handled events
			if msg.Type() == noise.SelfListening {
				log.Printf("Listening on: %s \n", msg.Payload())
				cancel() // stop listening for events
			}
		}
	}()

	// ... some code here
	// node.Dial("192.168.1.1:4008")
	// node.Close()

	// ... more code here
	node.Listen("127.0.0.1:4008")

}
