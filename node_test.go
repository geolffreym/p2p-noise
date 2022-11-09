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

func TestHandshake(t *testing.T) {
	out := new(bytes.Buffer)
	fl := log.Flags()
	log.SetFlags(0)
	log.SetOutput(out)

	nodeASocket := "127.0.0.1:9090"
	nodeBSocket := "127.0.0.1:9091"
	configurationA := config.New()
	configurationB := config.New()

	configurationA.Write(config.SetSelfListeningAddress(nodeASocket))
	configurationB.Write(config.SetSelfListeningAddress(nodeBSocket))

	t.Run("handshake A<->B", func(t *testing.T) {
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
