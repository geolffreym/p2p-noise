package noise

import (
	"bytes"
	"log"
	"os"
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

func TestTwoNodesHandshake(t *testing.T) {
	out := new(bytes.Buffer)
	fl := log.Flags()
	log.SetFlags(0)
	log.SetOutput(out)

	t.Run("handshake A<->B", func(t *testing.T) {
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

		var signals <-chan Signal = nodeA.Signals()
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
	log.Print(out)
}

func TestSomeNodesHandshake(t *testing.T) {
	out := new(bytes.Buffer)
	fl := log.Flags()
	log.SetFlags(0)
	log.SetOutput(out)

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

		var signalsA <-chan Signal = nodeA.Signals()
		for signalA := range signalsA {
			if signalA.Type() == NewPeerDetected {
				// Wait until new peer detected
				break
			}
		}

		var signalsB <-chan Signal = nodeB.Signals()
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

	log.SetFlags(fl)
	log.SetOutput(os.Stderr)
	log.Print(out)

}
