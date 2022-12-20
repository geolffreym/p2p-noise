package noise

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/geolffreym/p2p-noise/config"
)

// TODO benchmark memory usage for handshake byte pool
// TODO benchmark time for message exchange
// phase 1: adaptative lookup statwse
// phase 2: compression using brotli vs gzip
// phase 2 discovery module

func matchExpectedLogs(expectedBehavior []string, t *testing.T, f func()) {
	// store logs in buffer while the function run.
	out := new(bytes.Buffer)
	log.SetFlags(0)
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

func whenReadyForIncomingDial(nodes []*Node) *sync.WaitGroup {
	// Wait until all the nodes are ready for incoming connections.
	var wg sync.WaitGroup

	for _, node := range nodes {
		wg.Add(1)
		go node.Listen()
		// Populate wait group
		go func(n *Node) {
			signals, _ := n.Signals()
			for signal := range signals {
				if signal.Type() == SelfListening {
					wg.Done()
					return
				}
			}
		}(node)
	}

	return &wg
}

func TestWithZeroFutureDeadline(t *testing.T) {
	idle := futureDeadLine(0)

	if !idle.Equal(time.Time{}) {
		t.Errorf("Expected returned 'no deadline', got %v", idle)
	}

}

func TestNodesSecureMessageExchange(t *testing.T) {
	expected := "Hello node B"
	nodeASocket := "127.0.0.1:9090"
	nodeBSocket := "127.0.0.1:9091"
	configurationA := config.New()
	configurationB := config.New()

	var wg sync.WaitGroup
	wg.Add(1)

	configurationA.Write(config.SetSelfListeningAddress(nodeASocket))
	configurationB.Write(config.SetSelfListeningAddress(nodeBSocket))

	nodeA := New(configurationA)
	nodeB := New(configurationB)

	go nodeB.Listen()
	go nodeA.Listen()
	// Lets send a message from A to B and see if we receive the expected decrypted message
	go func(node *Node) {
		// Node A events channel
		signalsA, _ := node.Signals()
		for signalA := range signalsA {

			switch signalA.Type() {
			case SelfListening:
				wg.Done()
			case NewPeerDetected:
				// send a message to node b after handshake ready
				id := signalA.Payload() // here we receive the remote peer id
				_, err := nodeA.Send(id, []byte(expected))
				if err != nil {
					log.Print(err)
				}
			}
		}
	}(nodeA)

	wg.Wait()
	// Just dial to start handshake and close.
	nodeB.Dial(nodeASocket) // wait until handshake is done
	// Node B events channel
	signalsB, _ := nodeB.Signals()
	for signalB := range signalsB {
		if signalB.Type() == MessageReceived {
			// Wait until new peer detected
			got := signalB.Payload()
			if got != expected {
				t.Errorf("received message from node expected: %s, got %s", expected, got)
			}
			break
		}
	}

	// then just close nodes
	nodeA.Close()
	nodeB.Close()
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

		// wait until node is listening to start dialing
		whenReadyForIncomingDial([]*Node{nodeA, nodeB}).Wait()
		// Just dial to start handshake and close.
		nodeB.Dial(nodeASocket) // wait until handshake is done
		// then just close nodes
		nodeA.Close()
		nodeB.Close()
	})

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

	var nodes = []*Node{
		nodeA,
		nodeB,
		nodeC,
		nodeD,
	}

	// When all peers are listening then start dialing between them.
	whenReadyForIncomingDial(nodes).Wait()
	nodeB.Dial(nodeASocket)
	nodeC.Dial(nodeASocket)
	nodeC.Dial(nodeBSocket)
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

	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		b.StopTimer()

		var peers []*Node
		var peersNumber int = 10

		configurationA := config.New()
		configurationA.Write(config.SetSelfListeningAddress("127.0.0.1:"))
		nodeA := New(configurationA)

		for i := 0; i < peersNumber; i++ {
			address := "127.0.0.1:"
			configuration := config.New()
			configuration.Write(config.SetSelfListeningAddress(address))
			node := New(configuration)
			peers = append(peers, node)
		}

		fmt.Println("********************** Listen **********************")
		whenReadyForIncomingDial(append(peers, nodeA)).Wait()
		nodeAddress := nodeA.LocalAddr().String()

		// Start timer to measure the handshake process.
		// Handshake start when two nodes are connected and isn't happening before dial.
		// Avoid to add prev initialization.
		b.StartTimer()
		fmt.Println("********************** Dial **********************")
		for _, peer := range peers {
			peer.Dial(nodeAddress)
			peer.Close()
		}
	}

}
