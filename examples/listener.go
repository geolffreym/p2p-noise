//Copyright (c) 2022, Geolffrey Mena <gmjun2000@gmail.com>

//P2P Noise Secure handshake.
//
//See also: http://www.noiseprotocol.org/noise.html#introduction
package main

import (
	"context"
	"log"

	noise "github.com/geolffreym/p2p-noise"
	"github.com/geolffreym/p2p-noise/conf"
)

func main() {

	// Create settings from params and write in settings reference
	settings := conf.NewSettings()
	settings.Write(
		conf.SetMaxPeersConnected(10),
		conf.SetPeerDeadline(1800),
	)

	// Node factory
	node := noise.New(settings)
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
