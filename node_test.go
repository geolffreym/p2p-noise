package noise

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/geolffreym/p2p-noise/config"
	"golang.org/x/exp/slices"
)

func TestWithZeroFutureDeadline(t *testing.T) {
	idle := futureDeadLine(0)

	if !idle.Equal(time.Time{}) {
		t.Errorf("Expected returned 'no deadline', got %v", idle)
	}

}

func TestTwoNodesHandshake(t *testing.T) {
	out := new(bytes.Buffer)
	fl := log.Flags()
	log.SetFlags(0)
	log.SetOutput(out)

	expected_behavior := []string{
		"listening on 127.0.0.1:9091",     // node A listening
		"listening on 127.0.0.1:9090",     // node B listening
		"dialing to 127.0.0.1:9090",       // node B dialing to node A
		"starting handshake",              // Nodes starting handshake
		"generated ECDSA25519 public key", // Generating ECDSA Key Pair
		"generated X25519 public key",     //  Generating DH Key pair
		"sending e to remote",             // Handshake pattern
		"waiting for e from remote",
		"waiting for e, ee, s, es from remote",
		"sending e, ee, s, es to remote",
		"waiting for s, se from remote",
		"sending s, se to remote",
		"handshake complete", // Handshake complete
		"closing connections and shutting down node..",
	}

	t.Run("handshake A<->B trace", func(t *testing.T) {
		nodeASocket := "127.0.0.1:9090"
		nodeBSocket := "127.0.0.1:9091"
		configurationA := config.New()
		configurationB := config.New()

		configurationA.Write(config.SetSelfListeningAddress(nodeASocket))
		configurationB.Write(config.SetSelfListeningAddress(nodeBSocket))

		nodeA := New(configurationA)
		nodeB := New(configurationB)
		go nodeA.Listen()
		go nodeB.Listen()

		<-time.After(time.Second * 1)
		nodeB.Dial(nodeASocket)

		// Network events channel
		signals, _ := nodeA.Signals()
		for signal := range signals {
			if signal.Type() == NewPeerDetected {
				// Wait until new peer detected
				break
			}
		}

		nodeA.Close()
		nodeB.Close()
	})

	log.SetFlags(fl)
	log.SetOutput(os.Stderr)

	scanner := bufio.NewScanner(out)
	// The approach here is try to find the result in the expected behavior list.
	// If not found expected behavior in log results the test fail.
	for scanner.Scan() {
		got := scanner.Text()
		found := slices.Index(expected_behavior, got)
		// Not matched behavior
		if found < 0 {
			t.Errorf("expected to find '%s' behavior, got %d as not found", got, found)
		}
	}

}

func TestSomeNodesHandshake(t *testing.T) {
	t.Run("handshake N<->N", func(t *testing.T) {
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
		go nodeA.Listen()
		go nodeB.Listen()
		go nodeC.Listen()
		go nodeD.Listen()

		<-time.After(time.Second * 1)
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

		signalsB, _ := nodeB.Signals()
		for signalB := range signalsB {
			if signalB.Type() == NewPeerDetected {
				// Wait until new peer detected
				break
			}
		}

		nodeA.Close()
		nodeB.Close()
		nodeC.Close()
		nodeD.Close()
	})
}

// go test -benchmem -run=^$ -benchmem -memprofile memprofile.out -cpuprofile cpuprofile.out -bench=BenchmarkHandshakeProfile
// go tool pprof {file}
func BenchmarkHandshakeProfile(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b.ResetTimer()
		b.StopTimer()

		var peers []*Node
		var peersNumber int = 1
		nodeASocket := fmt.Sprintf("127.0.0.1:900%d", n)

		configurationA := config.New()
		configurationA.Write(config.SetSelfListeningAddress(nodeASocket))
		nodeA := New(configurationA)

		for i := 0; i < peersNumber; i++ {
			address := "127.0.0.1:"
			configuration := config.New()
			configuration.Write(config.SetSelfListeningAddress(address))
			node := New(configuration)
			peers = append(peers, node)
		}

		fmt.Println("********************** Listen **********************")
		go nodeA.Listen()
		for _, peer := range peers {
			go peer.Listen()
		}

		<-time.After(time.Second / 10)
		// Start timer to measure the handshake process.
		// Handshake start when two nodes are connected and isn't happening before dial.
		// Avoid to add prev initialization.
		b.StartTimer()
		fmt.Println("********************** Dial **********************")
		start := time.Now()
		for _, peer := range peers {
			peer.Dial(nodeASocket)
		}

		fmt.Printf("Took %v\n", time.Since(start))
	}
}

// TODO add test for message exchange encryption/decryption
