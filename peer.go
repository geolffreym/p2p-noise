package noise

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

// packet set needed properties to handle incoming message for peer.
type packet struct {
	// Ascending order for struct size
	Sig    []byte // N byte Signature
	Digest []byte // N byte Digest
}

// marshall encode packet to stream.
func marshall(p packet) bytes.Buffer {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	// encode packet to bytes
	encoder.Encode(p)
	return buffer
}

// unmarshall decode incoming message to packet.
func unmarshal(b []byte) packet {
	var p packet
	buf := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buf)
	// decode bytes to packet
	decoder.Decode(&p)
	return p
}

// peer its the trusty remote peer.
// Provide needed methods to interact with the secured session.
type peer struct {
	// Optimizing space with ordered types.
	// the attributes orders matters.
	// ref: https://stackoverflow.com/questions/2113751/sizeof-struct-in-go
	id   ID
	s    *session
	m    *metrics
	pool BytePool
}

// Create a new peer based on secure session
func newPeer(s *session) *peer {
	// Blake2 hashed remote public key.
	id := newBlake2ID(s.RemotePublicKey())
	return &peer{id, s, nil, nil}
}

// BindPool set a global memory pool for peer.
// Using pools remove latency from buffer allocation.
func (p *peer) BindPool(pool BytePool) {
	p.pool = pool
}

// Return peer id.
// Peer id its a blake2 hashed remote public key.
func (p *peer) ID() ID {
	return p.id
}

// Close its a forward method for internal `Close` method in session.
func (p *peer) Close() error {
	return p.s.Close()
}

// Close its a forward method for internal `SetDeadline` method in session.
func (p *peer) SetDeadline(t time.Time) error {
	return p.s.SetDeadline(t)
}

// Send send a message to Peer with size bundled in header for dynamic allocation of buffer.
// Each message is encrypted using session keys.
func (p *peer) Send(msg []byte) (uint32, error) {
	// Get a pool buffer chunk
	buffer := p.pool.Get()[:0]
	defer p.pool.Put(buffer)

	// Encrypt packet
	digest, err := p.s.Encrypt(buffer, msg)
	if err != nil {
		return 0, err
	}

	// Create a new packet to send it over the network
	sig := p.s.Sign(digest) // message signature
	packet := marshall(packet{sig, digest})
	bytes, err := p.s.Write(packet.Bytes())
	if err != nil {
		return 0, err
	}

	return uint32(bytes), nil
}

// Listen wait for incoming messages from Peer.
// Use the needed pool buffer based on incoming header.
func (p *peer) Listen() ([]byte, error) {
	// Get a pool buffer chunk
	buffer := p.pool.Get()
	defer p.pool.Put(buffer)

	bytes, err := p.s.Read(buffer)
	log.Printf("got %d bytes from peer", bytes)

	if err == nil {
		// decode incoming message
		inp := unmarshal(buffer)
		// validate message signature
		if !p.s.Verify(inp.Digest, inp.Sig) {
			err := fmt.Errorf("invalid signature for incoming message: %s", inp.Sig)
			return nil, errVerifyingSignature(err)
		}

		// Receive secure message from peer.
		// buffer[:0] means reset pivot to empty slice byte pool.
		return p.s.Decrypt(buffer[:0], inp.Digest)
	}

	// net: don't return io.EOF from zero byte reads
	// if err == io.EOF then peer connection is closed
	_, isNetError := err.(*net.OpError)
	if err != io.EOF && !isNetError {
		// end of message, but peer is still connected
		return nil, nil
	}

	// Close disconnected peer
	if err := p.s.Close(); err != nil {
		return nil, err
	}

	// Peer disconnected
	return nil, err

}
