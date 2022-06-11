//Copyright (c) 2022, Geolffrey Mena <gmjun2000@gmail.com>

//P2P Noise Secure handshake.
//
//See also: http://www.noiseprotocol.org/noise.html#introduction
package main

import (
	"context"
	"log"

	noise "github.com/geolffreym/p2p-noise"
)

func main() {
	node := noise.NewNode()
	// Network events channel
	ctx, cancel := context.WithCancel(context.Background())
	events := node.Events(ctx)

	go func() {
		for msg := range events {
			log.Printf("Listening on: %s \n", msg.Payload())
			cancel() // stop listening for events
		}
	}()

	// ... some code here
	// node.Dial("192.168.1.1:4008")
	// node.Close()

	// ... more code here
	node.Listen("127.0.0.1:4008")

}
