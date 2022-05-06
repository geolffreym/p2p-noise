package main

import (
	"fmt"

	"github.com/geolffreym/p2p-noise/node"
	"github.com/geolffreym/p2p-noise/pubsub"
)

func main() {

	listenAddr := "127.0.0.1:4007"
	listenAddrB := "127.0.0.1:4008"

	nodeA := node.New()
	nodeB := node.New()

	nodeA.Observe(func(msg *pubsub.Message) bool {
		switch msg.Type {
		case pubsub.SELF_LISTENING:
			fmt.Printf("Listening on: %s", msg.Payload)
		case pubsub.NEWPEER_DETECTED:
			fmt.Printf("New peer: %s", msg.Payload)
		case pubsub.MESSAGE_RECEIVED:
			fmt.Printf("New message: %s", msg.Payload)
		default:

		}

		return true
	})

	nodeA.Listen(listenAddr)

	nodeB.Listen(listenAddrB)
	nodeB.Observe(func(msg *pubsub.Message) bool {
		switch msg.Type {
		case pubsub.SELF_LISTENING:
			fmt.Printf("Listening on: %s", msg.Payload)
		case pubsub.NEWPEER_DETECTED:
			fmt.Printf("New peer: %s", msg.Payload)
		case pubsub.MESSAGE_RECEIVED:
			fmt.Printf("New message: %s", msg.Payload)
		}

		return true
	})

	nodeB.Dial(listenAddr)

	<-nodeA.Done

	// // Type assertion.. is b string type?
	// var b interface{} = "hello"
	// a, ok := b.(string)

	// fmt.Print(a)
	// fmt.Print(ok)

}
