# P2P Noise

[![Go](https://github.com/geolffreym/p2p-noise/actions/workflows/go.yml/badge.svg)](https://github.com/geolffreym/p2p-noise/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/geolffreym/p2p-noise.svg)](https://pkg.go.dev/github.com/geolffreym/p2p-noise)
[![Go Report Card](https://goreportcard.com/badge/github.com/geolffreym/p2p-noise)](https://goreportcard.com/report/github.com/geolffreym/p2p-noise)
[![codecov](https://codecov.io/gh/geolffreym/p2p-noise/branch/main/graph/badge.svg?token=TAI49WYVTS)](https://codecov.io/gh/geolffreym/p2p-noise)

P2P Noise library aims to serve as a tool to create secure P2P networks based on the Noise Framework.

* Quick creation of custom P2P networks.
* Small and easy to use.
* Simplistic and lightweight.

## Features

> [Noise Secure Handshake](http://www.noiseprotocol.org/):
Noise is a framework for building crypto protocols. Noise protocols support mutual and optional authentication, identity hiding, forward secrecy, zero round-trip encryption, and other advanced features.

> [Adaptive Lookup for Unstructured Peer-to-Peer Overlays](https://arxiv.org/pdf/1509.04417.pdf):
Most of the unstructured peer-to-peer overlays do not provide any performance guarantee. "Adaptive Lookup" propose a novel Quality of Service enabled lookup for unstructured peer-to-peer overlays that will allow the user’s query to traverse only those overlay links which satisfy the given constraints.

## Install

```
go get github.com/geolffreym/p2p-noise
```

## Basic usage

```
import (
	"context"
	"log"

	noise "github.com/geolffreym/p2p-noise"
	"github.com/geolffreym/p2p-noise/config"
)

func main() {

	// Create configuration from params and write in configuration reference
	configuration := config.New()
	configuration.Write(
		config.SetMaxPeersConnected(10),
		config.SetPeerDeadline(1800),
	)

	// Node factory
	node := noise.New(configuration)
	// Network events channel
	ctx, cancel := context.WithCancel(context.Background())
	var signals <-chan noise.SignalCtx = node.Signals(ctx)

	go func() {
		for signal := range signals {
			// Here could be handled events
			switch signal.Type() {
			// When a new peer is connected. Start ping pong game.
			case noise.NewPeerDetected:
				log.Printf("New Peer connected: %s \n", signal.Payload())
				signal.Reply([]byte("ping")) // start game
			// When we receive a message, check the content message and reply "ping" or "pong"
			case noise.MessageReceived:
				message := string(signal.Payload())
				if message == "ping" {
					signal.Reply([]byte("pong"))
				} else {
					signal.Reply([]byte("ping"))
				}
			// What we do when a peer get disconnected?
			case noise.PeerDisconnected:
				log.Printf("Peer disconnected")
				cancel() // stop listening for events
			}
		}
	}()

	// ... some code here
	err := node.Dial("192.168.1.17:8010")
	if err != nil {
		log.Fatal(err)
	}
	// node.Close()

	// ... more code here
	node.Listen()

}

```

## Development

Some available capabilities for dev support:

* **Run Tests**: `make test`

* **Build**: `make build`

* **Test Coverage**: `make coverage`

* **Benchmark**: `make benchmark`

* **Profiling**: `make profiling`

* **Code check**: `make code-check`

* **Code format**: `make code-fmt`

* **Flush cache**: `make clean`

* **Code Analysis**: `make check`

* **Build**: `make build`

Note: Please check [Makefile](https://github.com/geolffreym/p2p-noise) for more capabilities.  

## More info

* [Examples](https://github.com/geolffreym/p2p-noise) directory contains advanced examples of usage.
* For help or bugs please [create an issue](https://github.com/geolffreym/p2p-noise/issues).

*Special Thanks to @aphelionz.*
