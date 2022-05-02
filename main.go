package main

import (
	"fmt"
	"log"

	"github.com/geolffreym/p2p-network/network"
	"github.com/geolffreym/p2p-network/node"
)

func main() {

	// TODO Bootstrap node
	// TODO chat connected with bootstrap private network
	// const loopBackAddress string = "127.0.0.1"
	// const loopBackPort uint16 = 4002
	// node := networking.New(loopBackAddress, loopBackPort)
	// node.Listen().RunHandshake() // Run node and wait for incoming messages

	// net.Dial("tcp", fmt.Sprintf("%s:%d", loopBackAddress, loopBackPort))

	// a := runtime.GOMAXPROCS(0)
	// fmt.Print(a)
	const A = "127.0.0.1:4007"
	const B = "127.0.0.1:4008"

	nodeA, err := node.New(A).Listen()
	nodeA.Network.SetHandler(func(route *network.Route) {
		for {
			buf := make([]byte, 1024)
			_, err := route.Stream().Read(buf)
			if err == nil {
				fmt.Printf("Receiving new message: %s", buf)
			}

			fmt.Printf("%s", err)

			fmt.Printf("Sending pong to : %s", route.Socket())
			route.Stream().Write([]byte("Pong"))
		}

	})

	if err != nil {
		log.Fatal(err)
	}

	nodeB, err := node.New(B).Listen()
	nodeB.Network.SetHandler(func(route *network.Route) {
		for {
			// how i want to handle streaming between nodes?
			/*
				type Route struct {
					socket socket
					conn   net.Conn
					close  chan bool
				}
			*/
			buf := make([]byte, 1024)
			_, err := route.Stream().Read(buf)
			if err == nil {
				fmt.Printf("Receiving new message: %s", buf)
			}

			fmt.Printf("Sending ping\n to %s", route.Socket())
			route.Connection().Write([]byte("Ping"))
		}
	})

	// nodeB.S
	if err != nil {
		log.Fatal(err)
	}

	_, err = nodeB.Dial(A)
	// nodeB.S
	if err != nil {
		log.Fatal(err)
	}

	<-nodeA.Done

	// Type assertion.. is b string type?
	var b interface{} = "hello"
	a, ok := b.(string)

	fmt.Print(a)
	fmt.Print(ok)

}
