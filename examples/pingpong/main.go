package main

import (
	"context"
	"flag"
	"log"

	noise "github.com/geolffreym/p2p-noise"
	"github.com/geolffreym/p2p-noise/config"
)

var initiator bool
var ip, port string

func init() {
	// Node factory
	flag.BoolVar(&initiator, "i", false, "I start the game")
	flag.StringVar(&ip, "ip", "127.0.0.1", "IP address to connect")
	flag.StringVar(&port, "port", "8010", "Port to connect")
}

func main() {

	// Create configuration from params and write in configuration reference
	configuration := config.New()
	configuration.Write(
		config.SetMaxPeersConnected(10),
		config.SetPeerDeadline(1800),
	)

	// parse cli params
	flag.Parse()
	remote := ip + ":" + port
	node := noise.New(configuration)

	// Network events channel
	ctx, cancel := context.WithCancel(context.Background())
	var signals <-chan noise.Signal = node.Signals(ctx)

	go func() {
		// Wait for incoming message channel.
		for signal := range signals {
			// Here could be handled events
			switch signal.Type() {
			case noise.NewPeerDetected:
				// When a new peer is connected. Start ping pong game.
				log.Printf("New Peer connected: %x \n", signal.Payload())
				signal.Reply([]byte("ping")) // start game
				// TODO exchange peers?
				// TODO discovery module in action here?

			case noise.MessageReceived:
				// When we receive a message, check the content message and reply "ping" or "pong"
				message := string(signal.Payload())
				log.Printf("New Message: %s", message)
				if message == "ping" {
					signal.Reply([]byte("pong"))
				} else {
					signal.Reply([]byte("ping"))
				}

			case noise.PeerDisconnected:
				// What we do when a peer get disconnected?
				log.Printf("Peer disconnected")
				cancel() // stop listening for events
			}
		}
	}()

	// If i start the game then i should start dialing
	// If i am the second player i should just wait :)
	if initiator {
		// ... some code here
		log.Printf("dialing to %s", remote)
		err := node.Dial(remote)
		if err != nil {
			log.Fatal(err)
		}
	}

	// ... more code here
	node.Listen()

}
