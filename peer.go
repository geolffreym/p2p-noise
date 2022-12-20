package noise

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

// packet set needed properties to handle incoming message for peer.
type packet struct {
	// Ascending order for struct size
	Len int    // 8 bytes. Size of message
	Sig string // 16 byte Signature
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
	// Encrypt packet
	digest, err := p.s.Encrypt(msg)
	if err != nil {
		return 0, err
	}

	size := len(digest)     // the msg size
	sig := p.s.Sign(digest) // message signature

	// Create a new packet to send it over the network
	packet := packet{size, string(sig)}
	err = binary.Write(p.s, binary.BigEndian, packet)
	if err != nil {
		return 0, err
	}

	// Send encrypted message
	bytes, err := p.s.Write(digest)
	if err != nil {
		return 0, err
	}

	return uint32(bytes), nil
}

// Listen wait for incoming messages from Peer.
// Use the needed pool buffer based on incoming header.
func (p *peer) Listen() ([]byte, error) {
	var inp packet // incoming packet
	err := binary.Read(p.s, binary.BigEndian, &inp)
	if err != nil {
		return nil, err
	}

	// Get a pool buffer chunk
	buffer := p.pool.Get()[:inp.Len]
	defer p.pool.Put(buffer)

	// Sync buffered IO reading
	if _, err = p.s.Read(buffer); err == nil {
		// validate message signature
		if !p.s.Verify(buffer, []byte(inp.Sig)) {
			err := fmt.Errorf("invalid signature for incoming message: %s", inp.Sig)
			return nil, errVerifyingSignature(err)
		}
		// Receive secure message from peer.
		return p.s.Decrypt(buffer)
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
