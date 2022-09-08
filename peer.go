package noise

import (
	"encoding/binary"
	"encoding/hex"
	"io"
	"log"
	"net"
	"unsafe"

	"golang.org/x/crypto/blake2b"
)

// [ID] it's identity provider for peer.
type ID string

// Bytes return a byte slice representation for id.
func (i ID) Bytes() []byte {
	// A pointer value can't be converted to an arbitrary pointer type.
	// ref: https://go101.org/article/unsafe.html
	// no-copy conversion
	// ref: https://github.com/golang/go/issues/25484
	return *(*[]byte)(unsafe.Pointer(&i))
}

// Hash return a string representation for blake2 hash.
func (i ID) Hash() string {
	hash, err := blake2b.New(blake2b.Size256, nil)
	if err != nil {
		return ""
	}

	hash.Write(i.Bytes())
	digest := hash.Sum(nil)
	return hex.EncodeToString(digest)
}

// String return a string representation for 32-bytes hash.
func (i ID) String() string {
	return (string)(i)

}

// msgHeader set needed properties to handle incoming message for peer.
// Optimizing space with ordered types. Descending order.
// ref: https://stackoverflow.com/questions/2113751/sizeof-struct-in-go
type msgHeader struct {
	Hash  string // 16 bytes. string alias ID
	Len   uint32 // 4 bytes. Size of message
	Nonce uint32 // 4 bytes. Current message nonce
	Type  uint8  // 1 bytes. it's a number to handle message type.
}

// peer extends [net.Conn] interface.
// Each peer keep needed methods to interact with it.
// Please see [Connection Interface] for more details.
//
// [Connection Interface]: https://pkg.go.dev/net#Conn
type peer struct {
	net.Conn // TODO secure connection here?
	nonce    uint32
}

func newPeer(conn net.Conn) *peer {
	// Go does not provide the typical, type-driven notion of sub-classing,
	// but it does have the ability to “borrow” pieces of an implementation by embedding types within a struct or interface.
	return &peer{conn, 0}
}

// Return peer blake2 hash.
func (p *peer) ID() ID {
	// Temporary seed for ID. here could be used MAC, public key, etc..
	seed := p.RemoteAddr().String()
	return ID(seed)
}

// Send send a message to Peer with size bundled in header for dynamic allocation of buffer.
func (p *peer) Send(msg []byte) (int, error) {
	// TODO send msgHeader here
	// TODO add nonce ordered number to header
	// TODO Encrypt here with local key
	// write 4-bytes size header to share payload size
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
	// TODO decrypt here with remote key
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
