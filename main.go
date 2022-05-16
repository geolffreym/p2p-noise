//Copyright (c) 2022, Geolffrey Mena <gmjun2000@gmail.com>

//P2P Noise Secure handshake.
//
//See also: http://www.noiseprotocol.org/noise.html#introduction
package main

import (
	"fmt"

	"github.com/geolffreym/p2p-noise/network"
	"github.com/geolffreym/p2p-noise/node"
)

func main() {

	listenAddr := "127.0.0.1:4007"
	listenAddrB := "127.0.0.1:4008"
	listenAddrC := "127.0.0.1:4009"

	nodeA := node.NewNode()
	nodeB := node.NewNode()
	nodeC := node.NewNode()

	nodeA.Observe(func(msg *network.Message) bool {
		switch msg.Type {
		case network.SELF_LISTENING:
			fmt.Printf("Listening A on: %s \n", msg.Payload)
		case network.NEWPEER_DETECTED:
			fmt.Printf("New peer A: %s \n", msg.Payload)
		case network.MESSAGE_RECEIVED:
			fmt.Printf("New message A: %s \n", msg.Payload)
			msg.Reply([]byte("Pong"))

		default:

		}

		return true
	})

	nodeB.Observe(func(msg *network.Message) bool {

		switch msg.Type {
		case network.SELF_LISTENING:
			fmt.Printf("Listening B on: %s \n", msg.Payload)
		case network.NEWPEER_DETECTED:
			fmt.Printf("New peer B: %s \n", msg.Payload)
		case network.MESSAGE_RECEIVED:
			fmt.Printf("New message B: %s \n", msg.Payload)
			msg.Reply([]byte("Ping"))

		case network.CLOSED_CONNECTION:
			fmt.Print("Closed connection:")
		default:
		}

		return true
	})

	nodeC.Observe(func(msg *network.Message) bool {

		switch msg.Type {
		case network.SELF_LISTENING:
			fmt.Printf("Listening C on: %s \n", msg.Payload)
		case network.NEWPEER_DETECTED:
			fmt.Printf("New peer C: %s \n", msg.Payload)
		case network.MESSAGE_RECEIVED:
			fmt.Printf("New message C: %s \n", msg.Payload)
			msg.Peer.Send([]byte("Pong"))

		default:
		}

		return true
	})

	nodeA.Listen(listenAddr)
	nodeB.Listen(listenAddrB)
	nodeC.Listen(listenAddrC)

	nodeB.Dial(listenAddr)
	nodeB.Dial(listenAddrC)
	// time.Sleep(5 * time.Second)

	nodeB.Unicast(network.Socket(listenAddrC), []byte("Ping"))
	nodeB.Unicast(network.Socket(listenAddr), []byte("Ping"))
	// time.Sleep(5 * time.Second)
	// nodeA.Close()

	<-nodeB.Sentinel

	// // Type assertion.. is b string type?
	// var b interface{} = "hello"
	// a, ok := b.(string)

	// fmt.Print(a)
	// fmt.Print(ok)

}
