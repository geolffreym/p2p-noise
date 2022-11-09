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

	nodeB := New(configurationB)

	// t.Run(fmt.Sprintf("%x", e), func(t *testing.T) {
	// 	// Match recently added peer
	// 	if _, ok := router.Table()[e]; !ok {
	// 		t.Errorf("expected routed socket %#v", e.String())
	// 	}
	// })

	go func() {
		nodeA := New(configurationA)
		go nodeA.Listen()

		var signals <-chan Signal = nodeA.Signals(ctx)
		for signal := range signals {
			if signal.Type() == NewPeerDetected {
				log.Printf("%x", signal.Payload())
				<-time.After(time.Second * 5)
				nodeA.Close()
				close()
				return
			}
		}

	}()

	<-time.After(time.Second * 1)
	nodeB.Dial(nodeASocket)
	nodeB.Listen()

}
