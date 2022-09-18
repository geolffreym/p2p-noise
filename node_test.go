package noise

// func TestHandshake(t *testing.T) {
// 	nodeASocket := "127.0.0.1:9090"
// 	nodeBSocket := "127.0.0.1:9091"
// 	configurationA := config.New()
// 	configurationB := config.New()

// 	ctx, cancel := context.WithCancel(context.Background())
// 	configurationA.Write(config.SetSelfListeningAddress(nodeASocket))
// 	configurationB.Write(config.SetSelfListeningAddress(nodeBSocket))

// 	nodeA := New(configurationA)
// 	nodeB := New(configurationB)

// 	go func() {
// 		go func(n *Node) {
// 			var signals <-chan Signal = nodeA.Signals(ctx)
// 			for signal := range signals {
// 				if signal.Type() == NewPeerDetected {
// 					log.Print("Closing")
// 					n.Close()
// 					cancel()
// 				}
// 			}
// 		}(nodeA)
// 		nodeA.Listen()
// 	}()

// 	<-time.After(time.Second * 1)
// 	nodeB.Dial(nodeASocket)
// 	nodeB.Listen()

// }
