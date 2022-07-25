//Copyright (c) 2022, Geolffrey Mena <gmjun2000@gmail.com>

//P2P Noise Secure handshake.
//
//See also: http://www.noiseprotocol.org/noise.html#introduction
package main

import (
	"context"

	noise "github.com/geolffreym/p2p-noise"
	"github.com/geolffreym/p2p-noise/config"
)

func handshake() {

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
				// TODO handle here handshake logic
				cancel() // stop listening for events
			}
		}
	}()

	// ... some code here
	// node.Dial("192.168.1.1:4008")
	// node.Close()

	// ... more code here
	node.Listen()

}
