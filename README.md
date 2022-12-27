# P2P Noise

[![Go](https://github.com/geolffreym/p2p-noise/actions/workflows/go.yml/badge.svg)](https://github.com/geolffreym/p2p-noise/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/geolffreym/p2p-noise.svg)](https://pkg.go.dev/github.com/geolffreym/p2p-noise)
[![Go Report Card](https://goreportcard.com/badge/github.com/geolffreym/p2p-noise)](https://goreportcard.com/report/github.com/geolffreym/p2p-noise)
[![codecov](https://codecov.io/gh/geolffreym/p2p-noise/branch/main/graph/badge.svg?token=TAI49WYVTS)](https://codecov.io/gh/geolffreym/p2p-noise)

P2P Noise library aims to serve as a tool to create secure P2P networks based on the Noise Framework.

* Small and secure.
* Simplistic and lightweight.
* Quick creation of custom P2P networks.
* Modern Crypto Stack:
  * [Blake2 Hash](https://www.blake2.net/)
  * [ED255519 Signature](https://ed25519.cr.yp.to/)
  * [ChaCha20-Poly1305 Cypher]( https://en.wikipedia.org/wiki/ChaCha20-Poly1305)
  * [Diffie-Hellman Curve25519 Key Exchange](https://en.wikipedia.org/wiki/Curve25519)

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

```package main

import (
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

```

## Benchmarking

### Handshake Benchmark

Using [perflock](https://github.com/aclements/perflock) to prevent our benchmarks from using too much CPU at once.

```text
perflock -governor=80% go test -benchmem -run=^$ -benchtime 1s -bench=BenchmarkHandshakeProfile -cpu 1,2,4,8 -count 2 -memprofile mem.prof -cpuprofile cpu.prof
goos: linux
goarch: amd64
pkg: github.com/geolffreym/p2p-noise
cpu: Intel(R) Xeon(R) CPU E3-1505M v5 @ 2.80GHz
BenchmarkHandshakeProfile            962           1061028 ns/op           34391 B/op        267 allocs/op
BenchmarkHandshakeProfile           1322            919647 ns/op           34383 B/op        267 allocs/op
BenchmarkHandshakeProfile-2         1294            834506 ns/op           37095 B/op        285 allocs/op
BenchmarkHandshakeProfile-2         1398            858845 ns/op           36872 B/op        284 allocs/op
BenchmarkHandshakeProfile-4         1395            875618 ns/op           41912 B/op        323 allocs/op
BenchmarkHandshakeProfile-4         1341            914046 ns/op           41858 B/op        323 allocs/op
BenchmarkHandshakeProfile-8         1276            879535 ns/op           42055 B/op        324 allocs/op
BenchmarkHandshakeProfile-8         1279            929125 ns/op           41812 B/op        323 allocs/op
PASS
ok      github.com/geolffreym/p2p-noise 13.201s

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

* Chat with us joining to our [matrix room](https://matrix.to/#/!XgrTEPPGsKCPvdtDeC:matrix.org?via=matrix.org).
* [Examples](https://github.com/geolffreym/p2p-noise/examples) directory contains advanced examples of usage.
* For help or bugs please [create an issue](https://github.com/geolffreym/p2p-noise/issues).
