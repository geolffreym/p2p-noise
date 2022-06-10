//Copyright (c) 2022, Geolffrey Mena <gmjun2000@gmail.com>

//P2P Noise Secure handshake.
//
//See also: http://www.noiseprotocol.org/noise.html#introduction
package main

import (
	"log"

	noise "github.com/geolffreym/p2p-noise"
)

func main() {

	// Bind node to events messenger
	node := noise.NewNode()
	// // Every time that a new event is dispatched by node the messenger will notify to listener
	node.Intercept(func(msg noise.Message) {
		switch msg.Type() {
		case noise.SELF_LISTENING:
			log.Printf("Listening on: %s \n", msg.Payload())
		case noise.NEWPEER_DETECTED:
			log.Printf("New peer: %s \n", msg.Payload())
		case noise.CLOSED_CONNECTION:
			log.Print("Closed connection:")
		case noise.MESSAGE_RECEIVED:
			log.Printf("New message: %s \n", msg.Payload())
		default:

		}
	})

	node.Listen("127.0.0.1:4008")

}
