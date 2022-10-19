package noise

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/flynn/noise"
	"github.com/oxtoacart/bpool"
	"golang.org/x/crypto/chacha20poly1305"
)

// A CipherState provides symmetric encryption and decryption after a successful handshake.
type CipherState = *noise.CipherState

// BytePool implements a leaky pool of []byte in the form of a bounded channel.
type BytePool = *bpool.BytePool

// A DHKey is a keypair used for Diffie-Hellman key agreement.
type DHKey = noise.DHKey

type HandshakeState interface {
	// WriteMessage appends a handshake message to out. The message will include the
	// optional payload if provided. If the handshake is completed by the call, two
	// CipherStates will be returned, one is used for encryption of messages to the
	// remote peer, the other is used for decryption of messages from the remote
	// peer. It is an error to call this method out of sync with the handshake
	// pattern.
	WriteMessage(out, payload []byte) ([]byte, CipherState, CipherState, error)
	// ReadMessage processes a received handshake message and appends the payload,
	// if any to out. If the handshake is completed by the call, two CipherStates
	// will be returned, one is used for encryption of messages to the remote peer,
	// the other is used for decryption of messages from the remote peer. It is an
	// error to call this method out of sync with the handshake pattern.
	ReadMessage(out, message []byte) ([]byte, CipherState, CipherState, error)
	// PeerStatic returns the static key provided by the remote peer during
	// a handshake. It is an error to call this method if a handshake message
	// containing a static key has not been read.
	PeerStatic() []byte
	// MessageIndex returns the current handshake message id
	MessageIndex() int
}

// Buffer pools
// If bPools >= 1 a new buffered pool is created.
// If bPools == 0 a new no-buffered pool is created
const bPools = 1

// BLAKE2 is a cryptographic hash function faster than MD5, SHA-1, SHA-2, and SHA-3.
// [Blake2]: https://www.blake2.net/

// Cipher algorithm.
// ChaCha20-Poly1305 usually offers better performance than the more prevalent AES-GCM algorithm on systems where the CPU(s)
// does not feature the AES-NI instruction set extension.[2] As a result, ChaCha20-Poly1305 is sometimes preferred over
// AES-GCM due to its similar levels of security and in certain use cases involving mobile devices, which mostly use ARM-based CPUs.
var CipherSuite = noise.NewCipherSuite(noise.DH25519, noise.CipherChaChaPoly, noise.HashBLAKE2s)

// Default Handshake "XX" noise pattern.
// Our approach its use a balanced "time/security" pattern.
// Please see [NoisePatternExplorer] for more details.
//
// [NoisePatternExplorer]: https://noiseexplorer.com/patterns/XX/
var HandshakePattern = noise.HandshakeXX

// GenerateKeypair generates a new keypair using random as a source of entropy.
// Please see [Docs] for more details.
//
// [Docs]: http://www.noiseprotocol.org/noise.html#dh-functions
func generateKeyPair() (DHKey, error) {
	// Diffie-Hellman key pair
	var err error
	var kp DHKey

	// TODO should i persist seed?
	// TODO rand.Reader store it and retrieve it to avoid change the pub key in every new handshake?
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
		log.Print(err)
		err = fmt.Errorf("error creating handshake state: %v", err)
		return nil, errDuringHandshake(err)
	}

	return hs, nil
}

// noise config factory
// A Config provides the details necessary to process a Noise handshake. It is never modified by this package, and can be reused.
func newHandshakeConfig(initiator bool, kp noise.DHKey) noise.Config {
	return noise.Config{
		CipherSuite:   CipherSuite,
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
	s  *session // ref: https://go.dev/doc/effective_go#embedding
	hs HandshakeState
	p  BytePool
	i  bool
}

// Create a new handshake handler using provided connection and role.
func newHandshake(conn net.Conn, initiator bool) (*handshake, error) {
	kp, err := generateKeyPair()
	if err != nil {
		return nil, err
	}

	log.Printf("Generated public key %x", kp.Public)
	// set handshake state as initiator?
	conf := newHandshakeConfig(initiator, kp)
	// A HandshakeState tracks the state of a Noise handshake
	state, err := newHandshakeState(conf)
	if err != nil {
		return nil, err
	}

	// Setup the max of size possible for tokens exchanged between peers.
	// 64(DH keys) + 16(static key encrypted) + 2(size) = pool size
	size := 2*noise.DH25519.DHLen() + 2*chacha20poly1305.Overhead + 2
	pool := bpool.NewBytePool(bPools, size) // N pool of 84 bytes
	// Start a new session
	session := newSession(conn, kp)
	return &handshake{session, state, pool, initiator}, nil
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
	return h.hs.MessageIndex() >= len(HandshakePattern.Messages)
}

// Valid check if handshake sync is valid.
// If the process finished but the keys are not exchange as expected return error.
func (h *handshake) Valid(enc, dec CipherState) error {
	// This cyphers will be `not nil` until exchange patterns finish
	if h.Finish() && (enc == nil || dec == nil) {
		err := errors.New("invalid enc/dec cipher after handshake")
		return errDuringHandshake(err)
	}

	return nil
}

func (h *handshake) Start() error {
	if h.i {
		// Run as initiator role
		return h.Initiate()
	}
	// Run as remote "peer" role
	return h.Answer()
}

// Initiate start a new handshake with peer as a "dialer".
func (h *handshake) Initiate() error {
	// Send initial #1 message
	// bytes size = DHLen for e = ephemeral key
	log.Print("Sending e to remote")
	enc, dec, err := h.Send()
	if err != nil {
		err = fmt.Errorf("error sending `e` state: %v", err)
		return errDuringHandshake(err)
	}

	// Receive message #2 stage
	log.Print("Waiting for e, ee, s, es from remote")
	enc, dec, err = h.Receive()
	if err != nil {
		err = fmt.Errorf("error receiving `e, ee, s, es` state: %v", err)
		return errDuringHandshake(err)
	}

	// Send last handshake message #3 stage
	log.Print("Sending s, se to remote")
	enc, dec, err = h.Send()
	if err != nil {
		err = fmt.Errorf("error sending `s, se` state: %v", err)
		return errDuringHandshake(err)
	}

	// Check if synchronization is valid
	if err := h.Valid(enc, dec); err != nil {
		return err
	}

	// Bound handshake state to session
	h.s.SetState(h.hs)
	// Add keys for encrypt/decrypt operations in session.
	h.s.SetCyphers(enc, dec)
	return nil

}

// Answer start an answer for remote peer handshake request.
func (h *handshake) Answer() error {
	// Receive message #1 stage
	log.Print("Waiting for e from remote")
	enc, dec, err := h.Receive()
	if err != nil {
		err = fmt.Errorf("error receiving `e` state: %v", err)
		return errDuringHandshake(err)
	}

	// Send answer message #2 stage
	log.Print("Sending e, ee, s, es to remote")
	enc, dec, err = h.Send()
	if err != nil {
		err = fmt.Errorf("error sending `e, ee, s, es` state: %v", err)
		return errDuringHandshake(err)
	}

	// Receive message #2 stage
	log.Print("Waiting for s, se from remote")
	enc, dec, err = h.Receive()
	if err != nil {
		err = fmt.Errorf("error receiving `s, se` state: %v", err)
		return errDuringHandshake(err)
	}

	// Check if synchronization is valid
	if err := h.Valid(enc, dec); err != nil {
		return err
	}

	// Bound handshake state to session
	h.s.SetState(h.hs)
	// Add keys for encrypt/decrypt operations in session.
	h.s.SetCyphers(enc, dec)
	return nil
}

// Send create a new token based on message pattern synchronization and send it to remote peer.
func (h *handshake) Send() (e, d CipherState, err error) {
	var msg []byte
	// Get a chunk of bytes from pool
	buffer := h.p.Get()
	defer h.p.Put(msg)

	// WriteMessage appends a handshake message to out. The message will include the
	// optional payload if provided. If the handshake is completed by the call, two
	// CipherStates will be returned, one is used for encryption of messages to the
	// remote peer, the other is used for decryption of messages from the remote
	// peer. It is an error to call this method out of sync with the handshake
	// pattern.
	msg, e, d, err = h.hs.WriteMessage(buffer[:0], nil)
	if err != nil {
		return
	}

	// 2 bytes of header size
	binary.Write(h.s, binary.BigEndian, uint16(len(msg)))
	if _, err = h.s.Write(msg); err != nil {
		return
	}

	return
}

// Receive get a token from remote peer and synchronize it with local peer handshake state.
func (h *handshake) Receive() (e, d CipherState, err error) {
	var size uint16 // read bytes size from header
	err = binary.Read(h.s, binary.BigEndian, &size)
	if err != nil {
		return
	}

	// With size sent get a chunk from pool
	// Maybe here we don't need a new pool, just getting a chunk of current could help?
	buffer := h.p.Get()[:size]
	defer h.p.Put(buffer)

	// Wait for incoming message from remote
	if _, err = h.s.Read(buffer); err != nil {
		return
	}

	// ReadMessage processes a received handshake message and appends the payload,
	// if any to out. If the handshake is completed by the call, two CipherStates
	// will be returned, one is used for encryption of messages to the remote peer,
	// the other is used for decryption of messages from the remote peer. It is an
	// error to call this method out of sync with the handshake pattern.
	_, e, d, err = h.hs.ReadMessage(nil, buffer)
	if err != nil {
		return
	}

	return
}
