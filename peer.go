package noise

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"time"
)

// TODO add msgHeader to each message
// msgHeader set needed properties to handle incoming message for peer.
// Optimizing space with ordered types. Descending order.
// ref: https://stackoverflow.com/questions/2113751/sizeof-struct-in-go
// type msgHeader struct {
// 	Hash  string // 16 bytes. string alias ID
// 	Len   uint32 // 4 bytes. Size of message
// 	Nonce uint32 // 4 bytes. Current message nonce
// 	Type  uint8  // 1 bytes. it's a number to handle message type.
// }

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
	hash := Blake2(plaintext)
	// Populate id
	copy(id[:], hash)
	return id
}

// peer its the trusty remote peer.
// Keep needed methods to interact with the secured session.
type peer struct {
	s *session
	p BytePool
	i ID
	n uint32
}

// TODO write here docs
func newPeer(s *session) *peer {
	// Blake2 hashed remote public key.
	id := newBlake2ID(s.State().PeerStatic())
	return &peer{s, nil, id, 0}
}

// BindPool set a global memory pool for peer.
// Using pools remove latency from buffer allocation.
func (p *peer) BindPool(pool BytePool) {
	p.p = pool
}

// Return peer id.
// Peer id its a blake2 hashed remote public key.
func (p *peer) ID() ID {
	return p.i
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
	digest, err := p.s.Encrypt(msg)
	if err != nil {
		return 0, err
	}

	// 4 bytes of size header
	size := uint32(len(digest))
	binary.Write(p.s, binary.BigEndian, size)
	// Send message to session encryption.
	bytes, err := p.s.Write(digest)
	if err != nil {
		return 0, err
	}

	return uint32(bytes + 4), nil
}

// Listen wait for incoming messages from Peer.
// Use the needed pool buffer based on incoming header.
func (p *peer) Listen(maxPayloadSize uint32) ([]byte, error) {
	var size uint32 // read bytes size from header
	err := binary.Read(p.s, binary.BigEndian, &size)
	if err != nil {
		return nil, err
	}

	if size > maxPayloadSize {
		log.Printf("max payload size exceeded: MaxPayloadSize = %d", maxPayloadSize)
		return nil, errExceededMaxPayloadSize(maxPayloadSize)
	}

	// Get a pool buffer chunk
	buffer := p.p.Get()[:size]
	defer p.p.Put(buffer)
	// Sync buffered IO reading
	if _, err = p.s.Read(buffer); err == nil {
		// Receive secure message from peer.
		return p.s.Decrypt(buffer)
	}

	// net: don't return io.EOF from zero byte reads
	// if err == io.EOF then peer connection is closed
	err, isNetError := err.(*net.OpError)
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
