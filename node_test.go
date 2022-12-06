package noise

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/geolffreym/p2p-noise/config"
)

func TestWithZeroFutureDeadline(t *testing.T) {
	idle := futureDeadLine(0)

	if !idle.Equal(time.Time{}) {
		t.Errorf("Expected returned 'no deadline', got %v", idle)
	}

}

func TestTwoNodesHandshakeTrace(t *testing.T) {
	out := new(bytes.Buffer)
	log.SetFlags(0)
	log.SetOutput(out)

	expectedBehavior := []string{
		"starting handshake", // Nodes starting handshake
		"handshake complete", // Handshake complete
	}

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
	nodeA.Close()
	nodeB.Close()

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
