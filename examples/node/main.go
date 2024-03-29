package main

import (
	noise "github.com/geolffreym/p2p-noise"
	"github.com/geolffreym/p2p-noise/config"
)

func main() {

	// Create configuration from params and write in configuration reference
	configuration := config.New()
	configuration.Write(
		config.SetMaxPeersConnected(10),
	)

	// Node factory
	node := noise.New(configuration)
	// Network events channel
	signals, cancel := node.Signals()

	go func() {
		for signal := range signals {
			// Here could be handled events
			if signal.Type() == noise.NewPeerDetected {
				cancel()
			}
		}
	}()

	// ... some code here
	// node.Dial("192.168.1.1:4008")
	// node.Close()

	// ... more code here
	node.Listen()

}
