package main

import (
	"fmt"

	"github.com/geolffreym/p2p-noise/network"
	"github.com/geolffreym/p2p-noise/node"
)

func main() {

	listenAddr := "127.0.0.1:4007"
	listenAddrB := "127.0.0.1:4008"

	nodeA := node.New()
	nodeB := node.New()

	//
	nodeA.AddListener(network.LISTENING, func(route *network.Peer, args ...any) {
		fmt.Printf("Listening on: %s", listenAddr)

	})

	nodeA.AddListener(network.NEWPEER, func(route *network.Peer, args ...any) {
		fmt.Printf("Peer connected: %s", route.Socket())

	})

	nodeA.AddListener(network.MESSAGE, func(route *network.Peer, args ...any) {
		message := args[0]
		fmt.Printf("New message: %s\n", message)
		route.Write([]byte("pong"))
	})

	nodeB.AddListener(network.LISTENING, func(route *network.Peer, args ...any) {
		fmt.Printf("Listening on: %s", listenAddrB)

	})

	nodeB.AddListener(network.NEWPEER, func(route *network.Peer, args ...any) {
		fmt.Printf("Peer connected: %s", route.Socket())
		route.Write([]byte("ping"))

	})

	nodeB.AddListener(network.MESSAGE, func(route *network.Peer, args ...any) {
		message := args[0]
		fmt.Printf("New message: %s\n", message)
		route.Write([]byte("ping"))
	})

	nodeA.Listen(listenAddr)
	nodeB.Listen(listenAddrB)
	nodeB.Dial(listenAddr)

	<-nodeA.Done

	// // Type assertion.. is b string type?
	// var b interface{} = "hello"
	// a, ok := b.(string)

	// fmt.Print(a)
	// fmt.Print(ok)

}
