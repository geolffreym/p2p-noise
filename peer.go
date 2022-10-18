package noise

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"time"

	"golang.org/x/crypto/blake2b"
)

// blake2 return a 32-bytes representation for blake2 hash.
func blake2(i []byte) []byte {
	// New returns a new hash.Hash computing the BLAKE2b checksum with a custom length.
	// A non-nil key turns the hash into a MAC. The key must be between zero and 64 bytes long.
	// The hash size can be a value between 1 and 64 but it is highly recommended to use
	// values equal or greater than:
	// - 32 if BLAKE2b is used as a hash function (The key is zero bytes long).
	// - 16 if BLAKE2b is used as a MAC function (The key is at least 16 bytes long).
	// When the key is nil, the returned hash.Hash implements BinaryMarshaler
	// and BinaryUnmarshaler for state (de)serialization as documented by hash.Hash.
	hash, err := blake2b.New(blake2b.Size256, nil)
	if err != nil {
		return nil
	}

	hash.Write(i)
	return hash.Sum(nil)
}

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

// msgHeader set needed properties to handle incoming message for peer.
// Optimizing space with ordered types. Descending order.
// ref: https://stackoverflow.com/questions/2113751/sizeof-struct-in-go
type msgHeader struct {
	ID    [32]byte // 32 bytes. ID
	Sig   [32]byte // 32 byte. Signature
	Len   uint32   // 4 bytes. Size of message
	Nonce uint32   // 4 bytes. Current message nonce
}

// peer its the trusty remote peer.
// Keep needed methods to interact with the secured session.
type peer struct {
	s     *session
	id    ID
	pool  BytePool
	nonce uint32 // nonce
	// handshakeAt: time.Now().String(),
	// 	peers:     []..,
}

// TODO write here docs
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
	digest, err := p.s.Encrypt(msg)
	if err != nil {
		return 0, err
	}

	// TODO add msgHeader to each message
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
	buffer := p.pool.Get()[:size]
	defer p.pool.Put(buffer)
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
