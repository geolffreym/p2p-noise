package noise

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/geolffreym/p2p-noise/config"
)

// phase 1: metrics for adaptive lookup
// phase 2: compression using brotli vs gzip
// phase 2 discovery module
func traceMessageBetweenTwoPeers(nodeB *Node, expected string) bool {
	ready := make(chan bool)

	// Node B events channel
	signalsB, cancel := nodeB.Signals()
	for signalB := range signalsB {
		if signalB.Type() == MessageReceived {
			// When a new message is received:
			// Underneath the message is verified with remote PublicKey and decrypted with DH SharedKey.
			got := signalB.Payload()
			cancel() // stop the signaling
			return got == expected

		}
	}

	close(ready)
	// by default expected not message received as expected
	return false
}

func matchExpectedLogs(expectedBehavior []string, t *testing.T, f func()) {
	// store logs in buffer while the function run.
	out := new(bytes.Buffer)
	log.SetFlags(0)
	// store log output in buffer
	log.SetOutput(out)

	f() // Exec code to get log snapshot

	// without reset log output = race condition
	log.SetFlags(log.Flags())
	log.SetOutput(os.Stderr)
	// Scan output logs.
	scanner := bufio.NewScanner(out)
	// The approach here is try to find the result in the expected behavior list.
	// If not found expected behavior in log results the test fail.
start:
	for _, expected := range expectedBehavior {

		// Resume scanner carriage in the last log and try to find the next expected
		for scanner.Scan() {
			got := scanner.Text()
			if got == expected {
				continue start
			}
		}

		if scanner.Err() == nil {
			// Not matched behavior
			t.Errorf("expected to find '%s' behavior", expected)
		}
	}

}

func whenReadyForIncomingDial(node *Node) <-chan bool {
	// Wait until all the nodes are ready for incoming connections.
	ready := make(chan bool)

	go node.Listen()
	// Populate wait group
	go func(n *Node) {
		signals, _ := n.Signals()
		for signal := range signals {
			if signal.Type() == SelfListening {
				ready <- true
				return
			}
		}
	}(node)

	return ready
}

func TestWithZeroFutureDeadline(t *testing.T) {
	idle := futureDeadLine(0)

	if !idle.Equal(time.Time{}) {
		t.Errorf("Expected returned 'no deadline', got %v", idle)
	}

}

func TestTwoNodesHandshakeTrace(t *testing.T) {

	expectedBehavior := []string{
		"starting handshake", // Nodes starting handshake
		"generated ECDSA25519 public key",
		"generated X25519 public key",
		"handshake complete", // Handshake complete
		"closing connections and shutting down node..",
	}

	// check if the log output match with expectedBehavior
	matchExpectedLogs(expectedBehavior, t, func() {
		nodeASocket := "127.0.0.1:9090"
		nodeBSocket := "127.0.0.1:9091"
		configurationA := config.New()
		configurationB := config.New()

		configurationA.Write(config.SetSelfListeningAddress(nodeASocket))
		configurationB.Write(config.SetSelfListeningAddress(nodeBSocket))

		nodeA := New(configurationA)
		nodeB := New(configurationB)

		go nodeA.Listen()
		// then just close nodes
		defer nodeA.Close()
		defer nodeB.Close()

		// wait until node is listening to start dialing
		<-whenReadyForIncomingDial(nodeA)

		// Just dial to start handshake and close.
		nodeB.Dial(nodeASocket) // wait until handshake is done

	})

}

func TestPoolBufferSizeForMessageExchange(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	configurationA := config.New()
	configurationB := config.New()

	byteSize := 1 << 4
	ready := make(chan bool)
	b := make([]byte, byteSize)

	rand.Read(b) // fill buffer with pseudorandom numbers
	expected := string(b)
	configurationA.Write(config.SetPoolBufferSize(byteSize))
	configurationB.Write(config.SetPoolBufferSize(byteSize))

	nodeA := New(configurationA)
	nodeB := New(configurationB)

	go nodeB.Listen()
	go nodeA.Listen()
	defer nodeA.Close()
	defer nodeB.Close()

	// Lets send a message from A to B and see
	// if we receive the expected decrypted message
	go func(node *Node) {
		// Node A events channel
		signalsA, _ := node.Signals()
		for signalA := range signalsA {
			switch signalA.Type() {
			case SelfListening:
				ready <- true
			case NewPeerDetected:
				// send a message to node b after handshake ready
				id := signalA.Payload() // here we receive the remote peer id
				// Start interaction with remote peer
				// Underneath the message is encrypted and signed with local Private Key before send.
				nodeA.Send(id, []byte(expected))
			}
		}
	}(nodeA)

	<-ready

	// Node B events channel
	nodeB.Dial(nodeA.LocalAddr().String())
	signalsB, cancel := nodeB.Signals()

	for signalB := range signalsB {
		if signalB.Type() == MessageReceived {
			// When a new message is received:
			// Underneath the message is verified with remote PublicKey and decrypted with DH SharedKey.
			got := signalB.Payload()
			cancel() // stop the signaling

			if got != expected {
				t.Errorf("expected valid message equal to %s", expected)
			}

		}
	}

}

func TestSomeNodesHandshake(t *testing.T) {

	nodeASocket := "127.0.0.1:9090"
	nodeBSocket := "127.0.0.1:9091"
	nodeCSocket := "127.0.0.1:9092"
	nodeDSocket := "127.0.0.1:9093"

	configurationA := config.New()
	configurationB := config.New()
	configurationC := config.New()
	configurationD := config.New()

	configurationA.Write(config.SetSelfListeningAddress(nodeASocket))
	configurationB.Write(config.SetSelfListeningAddress(nodeBSocket))
	configurationC.Write(config.SetSelfListeningAddress(nodeCSocket))
	configurationD.Write(config.SetSelfListeningAddress(nodeDSocket))

	nodeA := New(configurationA)
	nodeB := New(configurationB)
	nodeC := New(configurationC)
	nodeD := New(configurationD)

	// When all peers are listening then start dialing between them.
	<-whenReadyForIncomingDial(nodeA)

	nodeB.Dial(nodeASocket)
	nodeC.Dial(nodeASocket)
	nodeD.Dial(nodeBSocket)

	// Network events channel
	signalsA, _ := nodeA.Signals()
	for signalA := range signalsA {
		if signalA.Type() == NewPeerDetected {
			// Wait until new peer detected
			break
		}
	}

	nodeA.Close()
	nodeB.Close()
	nodeC.Close()
	nodeD.Close()

}

// go test -benchmem -run=^$ -benchmem -memprofile memprofile.out -cpuprofile cpuprofile.out -bench=BenchmarkHandshakeProfile
// go tool pprof {file}
func BenchmarkHandshakeProfile(b *testing.B) {

	// Discard logs to avoid extra allocations.
	log.SetOutput(ioutil.Discard)

	configurationA := config.New()
	configurationA.Write(
		config.SetPoolBufferSize(1 << 2),
	)

	nodeA := New(configurationA)
	go nodeA.Listen()
	defer nodeA.Close()

	<-whenReadyForIncomingDial(nodeA)

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {

		b.StopTimer()
		for pb.Next() {

			b.StopTimer()
			configuration := config.New()
			configuration.Write(
				config.SetPoolBufferSize(1 << 2),
			)

			node := New(configuration)

			// Start timer to measure the handshake process.
			// Handshake start when two nodes are connected and isn't happening before dial.
			// Avoid to add prev initialization.
			b.StartTimer()
			node.Dial(nodeA.LocalAddr().String())
			node.Close()

		}

	})

}

func BenchmarkNodesSecureMessageExchange(b *testing.B) {
	// Discard logs to avoid extra allocations.
	log.SetOutput(ioutil.Discard)

	ready := make(chan bool)
	configurationA := config.New()
	configurationB := config.New()

	nodeA := New(configurationA)
	go nodeA.Listen()
	defer nodeA.Close()

	// Lets send a message from A to B and see
	// if we receive the expected decrypted message
	go func(node *Node) {
		// Node A events channel
		signalsA, _ := node.Signals()
		for signalA := range signalsA {
			switch signalA.Type() {
			case SelfListening:
				ready <- true
			case MessageReceived:
				// When a new message is received:
				// Underneath the message is verified with remote PublicKey and decrypted with DH SharedKey.
				signalA.Reply([]byte("pong"))
			}
		}
	}(nodeA)

	// wait until node a gets ready
	<-ready

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {

		nodeB := New(configurationB)
		defer nodeB.Close()

		// wait until handshake is done
		nodeB.Dial(nodeA.LocalAddr().String())
		signalsB, cancel := nodeB.Signals()

		b.StartTimer()

		for pb.Next() {

			// we need to measure message exchange only so we start time here
			// sign + encryption + marshall + transmission
			// Node B events channel

			for signalB := range signalsB {
				switch signalB.Type() {
				case NewPeerDetected:
					// send a message to node b after handshake ready
					id := signalB.Payload() // here we receive the remote peer id
					// Start interaction with remote peer
					// Underneath the message is encrypted and signed with local Private Key before send.
					nodeB.Send(id, []byte("ping"))
				case MessageReceived:
					// When a new message is received:
					// Underneath the message is verified with remote PublicKey and decrypted with DH SharedKey.
					if signalB.Payload() == "pong" {
						cancel()
					}
				}
			}

		}

	})

}
