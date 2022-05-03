package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/geolffreym/p2p-noise/network"
	"github.com/geolffreym/p2p-noise/node"
)

func main() {

	port := flag.String("port", "4007", "listening port")
	dial := flag.String("dial", "", "dial addres")
	message := flag.String("message", "", "message to address")
	receptor := flag.String("receptor", "", "the receptor address for message")
	flag.Parse()

	node, err := node.New(*port).Listen()
	node.Network.AddEventListener(network.NEWPEER, func(route *network.Route, args ...any) {
		fmt.Printf("Peer connected: %s", route.Socket())

	})

	node.Network.AddEventListener(network.MESSAGE, func(route *network.Route, args ...any) {
		message := args[0]
		fmt.Printf("New message: %s\n", message)
	})

	if *dial != "" {
		_, err := node.Dial(*dial)
		fmt.Print(err)
	}

	if *message != "" && *receptor != "" {
		node.Unicast(network.Socket(*receptor), []byte(*message))
	}

	if err != nil {
		log.Fatal(err)
	}

	<-node.Done

	// // Type assertion.. is b string type?
	// var b interface{} = "hello"
	// a, ok := b.(string)

	// fmt.Print(a)
	// fmt.Print(ok)

}
