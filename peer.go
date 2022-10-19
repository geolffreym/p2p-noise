package noise

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

// [ID] it's identity provider for peer.
type ID [32]byte

// Bytes return a byte slice representation for id.
func (i ID) Bytes() []byte {
	return i[:]
}

// String return a string representation for 32-bytes hash.
func (i ID) String() string {
	return (string)(i[:])
}

// Create a new id blake2 hash based.
func newBlake2ID(plaintext []byte) ID {
	var id ID
	// Hash
	hash := blake2(plaintext)
	// Populate id
	copy(id[:], hash)
	return id
}

// packet set needed properties to handle incoming message for peer.
// Optimizing space with ordered types. Descending order.
// ref: https://stackoverflow.com/questions/2113751/sizeof-struct-in-go
type packet struct {
	// Ascending order for struct size
	Len uint32 // 4 bytes. Size of message
	Sig []byte // N byte Signature
}

// peer its the trusty remote peer.
// Keep needed methods to interact with the secured session.
type peer struct {
	s     *session
	id    ID
	pool  BytePool
	nonce uint32
	// handshakeAt: time.Now().String(),
	// 	peers:     []..,
}

// Create a new peer based on secure session
func newPeer(s *session) *peer {
	// Blake2 hashed remote public key.
	id := newBlake2ID(s.State().PeerStatic())
	return &peer{s, id, nil, 0}
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
	return p.SetDeadline(t)
}

// Send send a message to Peer with size bundled in header for dynamic allocation of buffer.
// Each message is encrypted using session keys.
func (p *peer) Send(msg []byte) (uint32, error) {
	// Encrypt packet
	digest, err := p.s.Encrypt(msg)
	if err != nil {
		return 0, err
	}

	size := uint32(len(digest)) // the msg size
	sig := p.s.Sign(digest)     // message signature

	// Create a new packet to send it over the network
	packet := packet{size, sig}
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
func (p *peer) Listen(maxPayloadSize uint32) ([]byte, error) {
	var inp packet // incoming packet
	err := binary.Read(p.s, binary.BigEndian, &inp)
	if err != nil {
		return nil, err
	}

	// If the size of the message in packet exceed may expected
	if inp.Len > maxPayloadSize {
		log.Printf("max payload size exceeded: MaxPayloadSize = %d", maxPayloadSize)
		return nil, errExceededMaxPayloadSize(maxPayloadSize)
	}

	// Get a pool buffer chunk
	buffer := p.pool.Get()[:inp.Len]
	defer p.pool.Put(buffer)

	// Sync buffered IO reading
	if _, err = p.s.Read(buffer); err == nil {
		// validate message signature
		if !p.s.Verify(buffer, inp.Sig) {
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
