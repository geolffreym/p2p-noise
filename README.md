# P2P Noise

[![Go](https://github.com/geolffreym/p2p-noise/actions/workflows/go.yml/badge.svg)](https://github.com/geolffreym/p2p-noise/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/geolffreym/p2p-noise.svg)](https://pkg.go.dev/github.com/geolffreym/p2p-noise)

P2P Noise library aims to serve as a tool to create secure P2P networks based on the Noise Framework.

* Quick creation of custom P2P networks
* Small and easy to use
* Simplistic and Lightweight

## Features

* [Noise](http://www.noiseprotocol.org/) Secure Handshake .
* [Adaptive Lookup for Unstructured
Peer-to-Peer Overlays](https://arxiv.org/pdf/1509.04417.pdf#:~:text=An%20unstructured%20P2P%20system%20is,(or%20scale%20free%20networks).)

## Install

```
go get github.com/geolffreym/p2p-noise
```

## Basic usage

```

node := noise.NewNode()
 // Network events channel
 ctx, cancel := context.WithCancel(context.Background())
 events := node.Events(ctx)

 go func() {
  for msg := range events {
   log.Printf("Listening on: %s \n", msg.Payload())
   cancel() // stop listening for events
  }
 }()

 // ... some code here
 // node.Dial("192.168.1.1:4008")
 // node.Close()

 // ... more code here
 node.Listen("127.0.0.1:4008")

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

* **Flush cache***: `make clean`

* **Code Analysis**: `make check`

* **Compile**: `make compile`

Note: `Compile` command will attempt to compile for every OS-arch, please check [MakeFile](https://github.com/geolffreym/p2p-noise) for more capabilities.  

## More info

* [Examples](https://github.com/geolffreym/p2p-noise) directory contain advanced examples of usage.
* For help or bugs please [create an issue](https://github.com/geolffreym/p2p-noise/issues).


*Special Thanks to @aphelionz for his patience and support.*
