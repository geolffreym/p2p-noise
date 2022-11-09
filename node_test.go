package noise

import (
	"context"
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

func TestHandshake(t *testing.T) {
	nodeASocket := "127.0.0.1:9090"
	nodeBSocket := "127.0.0.1:9091"
	configurationA := config.New()
	configurationB := config.New()

	ctx, close := context.WithCancel(context.Background())
	configurationA.Write(config.SetSelfListeningAddress(nodeASocket))
	configurationB.Write(config.SetSelfListeningAddress(nodeBSocket))
	nodeA := New(configurationA)
	nodeB := New(configurationB)

	go func(na *Node) {
		go func(n *Node) {
			var signals <-chan Signal = nodeA.Signals(ctx)
			for signal := range signals {
				if signal.Type() == NewPeerDetected {
					log.Printf("%x", signal.Payload())
					<-time.After(time.Second * 5)
					n.Close()
					close()
					return
				}
			}
		}(na)
		na.Listen()
	}(nodeA)

	<-time.After(time.Second * 1)
	nodeB.Dial(nodeASocket)
	nodeB.Listen()

}
