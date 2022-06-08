//Copyright (c) 2022, Geolffrey Mena <gmjun2000@gmail.com>

//P2P Noise Secure handshake.
//
//See also: http://www.noiseprotocol.org/noise.html#introduction
package main

import (
	"fmt"
	"time"

	"github.com/geolffreym/p2p-noise/network"
	"github.com/geolffreym/p2p-noise/node"
)

func main() {

	listenAddr := "127.0.0.1:4007"
	listenAddrB := "127.0.0.1:4008"

	nodeA := node.NewNode()
	nodeB := node.NewNode()

	nodeA.Observe(func(msg network.Message) {
		switch msg.Type() {
		case network.SELF_LISTENING:
			fmt.Printf("Listening A on: %s \n", msg.Payload())
		case network.NEWPEER_DETECTED:
			fmt.Printf("New peer A: %s \n", msg.Payload())
		case network.CLOSED_CONNECTION:
			fmt.Print("Closed connection A:")
		case network.MESSAGE_RECEIVED:
			fmt.Printf("New message A: %s \n", msg.Payload())

		default:

		}
	})

	nodeB.Observe(func(msg network.Message) {

		switch msg.Type() {
		case network.SELF_LISTENING:
			fmt.Printf("Listening B on: %s \n", msg.Payload())
		case network.NEWPEER_DETECTED:
			fmt.Printf("New peer B: %s \n", msg.Payload())
		case network.MESSAGE_RECEIVED:
			fmt.Printf("New message B: %s \n", msg.Payload())

		case network.CLOSED_CONNECTION:
			fmt.Print("Closed connection:")
		case network.PEER_DISCONNECTED:
			fmt.Printf("Peer disconnected: %s \n", msg.Payload())
		default:
		}

	})

	go nodeA.Listen(listenAddr)

	time.Sleep(1 * time.Second)
	nodeB.Dial(listenAddr)

	go func(node node.Node) {
		time.Sleep(5 * time.Second)
		node.Broadcast([]byte("Hola bebe"))
	}(nodeA)

	go func(node node.Node) {
		time.Sleep(10 * time.Second)
		node.Close()
	}(nodeA)

	nodeB.Listen(listenAddrB)
	// time.Sleep(5 * time.Second)
	// nodeA.Close()

	// // Type assertion.. is b string type?
	// var b interface{} = "hello"
	// a, ok := b.(string)

	// fmt.Print(a)
	// fmt.Print(ok)

}
