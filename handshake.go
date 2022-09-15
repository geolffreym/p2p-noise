package noise

import (
	"crypto/rand"
	"errors"
	"fmt"
	"hash"
	"log"
	"net"

	"github.com/flynn/noise"
	"github.com/oxtoacart/bpool"
	"golang.org/x/crypto/blake2b"
)

// BLAKE2 is a cryptographic hash function faster than MD5, SHA-1, SHA-2, and SHA-3.
// [Blake2]: https://www.blake2.net/
type blake2Fn struct{}

func (h blake2Fn) HashName() string { return "BLAKE2b" }
func (h blake2Fn) Hash() hash.Hash {
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
		panic(err)
	}

	return hash
}

// Cipher algorithm.
// ChaCha20-Poly1305 usually offers better performance than the more prevalent AES-GCM algorithm on systems where the CPU(s)
// does not feature the AES-NI instruction set extension.[2] As a result, ChaCha20-Poly1305 is sometimes preferred over
// AES-GCM due to its similar levels of security and in certain use cases involving mobile devices, which mostly use ARM-based CPUs.
var HashBLAKE2 noise.HashFunc = blake2Fn{}
var cipherSuite = noise.NewCipherSuite(noise.DH25519, noise.CipherChaChaPoly, HashBLAKE2)
var HandshakePattern = noise.HandshakeXX

// GenerateKeypair generates a new keypair using random as a source of entropy.
// Please see [Docs] for more details.
//
// [Docs]: http://www.noiseprotocol.org/noise.html#dh-functions
func generateKeyPair() (noise.DHKey, error) {
	// Diffie-Hellman key pair
	var err error
	var kp noise.DHKey

	kp, err = noise.DH25519.GenerateKeypair(rand.Reader)
	if err != nil {
		err := fmt.Errorf("error trying to generate `s` keypair: %w", err)
		return kp, errDuringHandshake(err)
	}

	return kp, nil
}

// NewHandshakeState starts a new handshake using the provided configuration.
// A HandshakeState tracks the state of a Noise handshake
// Please see [Handshake State] for more details.
//
// [Handshake State]: https://pkg.go.dev/github.com/flynn/noise#HandshakeState
func newHandshakeState(conf noise.Config) (*noise.HandshakeState, error) {
	// handshake state
	hs, err := noise.NewHandshakeState(conf)
	if err != nil {
		err = fmt.Errorf("error creating handshake state: %v", err)
		return nil, errDuringHandshake(err)
	}

	return hs, nil
}

// noise config factory
// A Config provides the details necessary to process a Noise handshake. It is never modified by this package, and can be reused.
func newHandshakeConfig(initiator bool, kp noise.DHKey) noise.Config {
	return noise.Config{
		CipherSuite:   cipherSuite,
		Pattern:       HandshakePattern,
		Initiator:     initiator,
		StaticKeypair: kp,
	}
}

// handshake execute the steps needed for the noise handshake XX pattern.
// Please see [XX Pattern] for more details. [XX Explorer] pattern.
//
// [XX Pattern]: http://www.noiseprotocol.org/noise.html#handshake-patterns
// [XX Explorer]: https://noiseexplorer.com/patterns/XX/
type handshake struct {
	conn  net.Conn // ref: https://go.dev/doc/effective_go#embedding
	state *noise.HandshakeState
	s     *session
}

func newHandshake(conn net.Conn, initiator bool) (*handshake, error) {
	kp, err := generateKeyPair()
	if err != nil {
		return nil, err
	}

	// set handshake state as initiator?
	conf := newHandshakeConfig(initiator, kp)
	// A HandshakeState tracks the state of a Noise handshake
	state, err := newHandshakeState(conf)
	if err != nil {
		return nil, err
	}

	return &handshake{
		conn, state,
		newSession(conn),
	}, nil
}

// Session return secured session after handshake.
// This session is invalid if handshake isn't finished.
func (h *handshake) Session() *session {
	if !h.Finish() {
		return nil
	}

	return h.s
}

// Finish return the handshake state.
// Return true if handshake is finished otherwise false.
func (h *handshake) Finish() bool {
	return h.state.MessageIndex() >= len(HandshakePattern.Messages)
}

// Valid check if handshake sync is valid.
// If the process finished but the keys are not exchange as expected return error.
func (h *handshake) Valid(enc, dec *noise.CipherState) error {
	// This cyphers will be `not nil` until exchange patterns finish
	if h.Finish() && (enc == nil || dec == nil) {
		err := errors.New("invalid enc/dec cipher after handshake")
		return errDuringHandshake(err)
	}

	return nil
}

func (h *handshake) Start(initiator bool) error {
	if initiator {
		// Run as initiator role
		return h.Initiate()
	}
	// Run as remote "peer" role
	return h.Answer()
}

// Initiate start a new handshake with peer as a "dialer".
func (h *handshake) Initiate() error {

	// Reserved pool buffer chunk
	// DHLEN = 32
	// For "s": Sets temp to the next DHLEN + 16 bytes of the message if HasKey() == True
	size := 2*noise.DH25519.DHLen() + 16
	pool := bpool.NewBytePool(size, size)
	buffer := pool.Get()
	defer pool.Put(buffer)

	// Send initial #1 message
	// bytes size = DHLen for e = ephemeral key
	log.Print("Sending e to remote")
	enc, dec, err := h.Send(buffer)
	if err != nil {
		err = fmt.Errorf("error sending `e` state: %v", err)
		return errDuringHandshake(err)
	}

	// Receive message #2 stage
	log.Print("Waiting for e, ee, s, es from remote")
	enc, dec, err = h.Receive(buffer)
	if err != nil {
		err = fmt.Errorf("error receiving `e, ee, s, es` state: %v", err)
		return errDuringHandshake(err)
	}

	// Send last handshake message #3 stage
	log.Print("Sending s, se to remote")
	enc, dec, err = h.Send(buffer)
	if err != nil {
		err = fmt.Errorf("error sending `s, se` state: %v", err)
		return errDuringHandshake(err)
	}

	// Check if synchronization is valid
	if err := h.Valid(enc, dec); err != nil {
		return err
	}

	// Add keys for encrypt/decrypt operations in session.
	h.s.SetCyphers(enc, dec)
	return nil

}

// Answer start an answer for remote peer handshake request.
func (h *handshake) Answer() error {
	// Reserved pool buffer chunk
	size := 2*noise.DH25519.DHLen() + 16
	pool := bpool.NewBytePool(size, size)
	buffer := pool.Get()
	defer pool.Put(buffer)

	// Receive message #1 stage
	log.Print("Waiting for e from remote")
	enc, dec, err := h.Receive(buffer)
	if err != nil {
		err = fmt.Errorf("error receiving `e` state: %v", err)
		return errDuringHandshake(err)
	}

	// Send answer message #2 stage
	log.Print("Sending e, ee, s, es to remote")
	enc, dec, err = h.Send(buffer)
	if err != nil {
		err = fmt.Errorf("error sending `e, ee, s, es` state: %v", err)
		return errDuringHandshake(err)
	}

	// Receive message #2 stage
	log.Print("Waiting for s, se from remote")
	enc, dec, err = h.Receive(buffer)
	if err != nil {
		err = fmt.Errorf("error receiving `s, se` state: %v", err)
		return errDuringHandshake(err)
	}

	// Check if synchronization is valid
	if err := h.Valid(enc, dec); err != nil {
		return err
	}

	// Add keys for encrypt/decrypt operations in session.
	h.s.SetCyphers(enc, dec)
	return nil
}

// Send create a new token based on message pattern synchronization and send it to remote peer.
func (h *handshake) Send(buffer []byte) (e, d *noise.CipherState, err error) {
	var msg []byte
	// WriteMessage appends a handshake message to out. The message will include the
	// optional payload if provided. If the handshake is completed by the call, two
	// CipherStates will be returned, one is used for encryption of messages to the
	// remote peer, the other is used for decryption of messages from the remote
	// peer. It is an error to call this method out of sync with the handshake
	// pattern.
	msg, e, d, err = h.state.WriteMessage(buffer, nil)
	if err != nil {
		return
	}

	// Send message to remote peer
	if _, err = h.conn.Write(msg); err != nil {
		return
	}

	return
}

// Receive get a token from remote peer and synchronize it with local peer handshake state.
func (h *handshake) Receive(buffer []byte) (e, d *noise.CipherState, err error) {
	// Wait for incoming message from remote
	if _, err = h.conn.Read(buffer); err != nil {
		return
	}

	// ReadMessage processes a received handshake message and appends the payload,
	// if any to out. If the handshake is completed by the call, two CipherStates
	// will be returned, one is used for encryption of messages to the remote peer,
	// the other is used for decryption of messages from the remote peer. It is an
	// error to call this method out of sync with the handshake pattern.
	_, e, d, err = h.state.ReadMessage(buffer, nil)
	if err != nil {
		return
	}

	return
}
