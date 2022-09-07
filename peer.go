package noise

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

// peer extends [net.Conn] interface.
// Each peer keep needed methods to interact with it.
// Please see [Connection Interface] for more details.
//
// [Connection Interface]: https://pkg.go.dev/net#Conn
type peer struct {
	net.Conn // embedded net.Conn to peer. ref: https://go.dev/doc/effective_go#embedding
	// TODO check: es realmente necesario mantener este socket aca?
	socket Socket // IP and Port address for peer. https://en.wikipedia.org/wiki/Network_socket
}

func newPeer(socket Socket, conn net.Conn) *peer {
	// Go does not provide the typical, type-driven notion of sub-classing,
	// but it does have the ability to “borrow” pieces of an implementation by embedding types within a struct or interface.
	return &peer{
		conn,
		socket,
	}
}

// Return peer socket.
// eg. "127.0.0.1:2000"
func (p *peer) Socket() Socket {
	// TODO check: que pasa si el address retornado por el peer es diferente al que se establece en el constructor?
	// TODO escribir una prueba en el router que verifique el indice del router sea el mismo socket retornado por el peer en este
	return Socket(p.RemoteAddr().String())
	// return p.socket
}

// Send send a message to Peer with size bundled in header for dynamic allocation of buffer.
func (p *peer) Send(msg []byte) (int, error) {
	// write 4-bytes size header to share payload size
	// TODO add origin peer id
	// TODO add type of message eg. handshake, literal..
	// TODO add nonce ordered number to header
	err := binary.Write(p, binary.BigEndian, uint32(len(msg)))
	if err != nil {
		return 0, err
	}

	// Write payload
	bytesSent, err := p.Write(msg)
	return bytesSent + 4, err
}

// Listen wait for incoming messages from Peer.
// Each message keep a header with message size to allocate buffer dynamically.
func (p *peer) Listen(maxPayloadSize uint32) ([]byte, error) {

	var size uint32 // read bytes size from header
	err := binary.Read(p, binary.BigEndian, &size)

	// Error trying to read `size`
	if err != nil {
		return nil, err
	}

	if size > maxPayloadSize {
		log.Printf("max payload size exceeded: MaxPayloadSize = %d", maxPayloadSize)
		return nil, errExceededMaxPayloadSize(maxPayloadSize)
	}

	// Dynamic allocation based on msg size
	buf := make([]byte, size)
	// Sync buffered IO reading
	if _, err = p.Read(buf); err == nil {
		// Sync incoming message
		return buf, nil
	}

	// net: don't return io.EOF from zero byte reads
	// if err == io.EOF then peer connection is closed
	err, isNetError := err.(*net.OpError)
	if err != io.EOF && !isNetError {
		// end of message, but peer is still connected
		return nil, nil
	}

	// Close disconnected peer
	if err := p.Close(); err != nil {
		return nil, err
	}

	// Peer disconnected
	return nil, err

}
