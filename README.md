# P2P Noise

[![Go](https://github.com/geolffreym/p2p-noise/actions/workflows/go.yml/badge.svg)](https://github.com/geolffreym/p2p-noise/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/geolffreym/p2p-noise.svg)](https://pkg.go.dev/github.com/geolffreym/p2p-noise)
[![Go Report Card](https://goreportcard.com/badge/github.com/geolffreym/p2p-noise)](https://goreportcard.com/report/github.com/geolffreym/p2p-noise)
[![codecov](https://codecov.io/gh/geolffreym/p2p-noise/branch/main/graph/badge.svg?token=TAI49WYVTS)](https://codecov.io/gh/geolffreym/p2p-noise)

P2P Noise library aims to serve as a tool to create secure P2P networks based on the Noise Framework.

* Quick creation of custom P2P networks.
* Simplistic and lightweight.
* Small and secure.

## Features

> [Blake2 Hashing](https://www.blake2.net/):
BLAKE2 is a cryptographic hash function faster than MD5, SHA-1, SHA-2, and SHA-3, yet is at least as secure as the latest standard SHA-3. BLAKE2 has been adopted by many projects due to its high speed, security, and simplicity.

> [End to End Encryption](https://en.wikipedia.org/wiki/End-to-end_encryption)
In E2EE, the data is encrypted on the sender's system or device, and only the intended recipient can decrypt it. As it travels to its destination, the message cannot be read or tampered with by an internet service provider (ISP), application service provider, hacker or any other entity or service.

> [ED255519 Signature](https://ed25519.cr.yp.to/):
A digital signature is a mathematical scheme for verifying the authenticity of digital messages or documents. A valid digital signature, where the prerequisites are satisfied, gives a recipient very high confidence that the message was created by a known sender (authenticity), and that the message was not altered in transit (integrity).

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
package main

import (
 "context"

 noise "github.com/geolffreym/p2p-noise"
 "github.com/geolffreym/p2p-noise/config"
)

func handshake() {

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
 var signals <-chan noise.Signal = node.Signals(ctx)

 go func() {
  for signal := range signals {
   // Here could be handled events
   if signal.Type() == noise.NewPeerDetected {
    cancel() // stop listening for events
   }
  }
 }()

 // ... some code here
 // node.Dial("192.168.1.1:4008")
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
* **Build**: `make build`

Note: Please check [Makefile](https://github.com/geolffreym/p2p-noise/Makefile) for more capabilities.  

## More info

* [Examples](https://github.com/geolffreym/p2p-noise) directory contains advanced examples of usage.
* For help or bugs please [create an issue](https://github.com/geolffreym/p2p-noise/issues).
