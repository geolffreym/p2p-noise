package noise

import (
	"crypto/ed25519"
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

// [DHKey] is a keypair used for Diffie-Hellman key agreement.
// Please see [docs] for more details.
//
// [docs]: http://www.noiseprotocol.org/noise.html#dh-functions
type DHKey = noise.DHKey
type PublicKey = ed25519.PublicKey
type PrivateKey = ed25519.PrivateKey

// [EDKeyPair] hold public/private using entropy from rand.
// Every new handshake generate a new key pair.
type EDKeyPair struct {
	Private PrivateKey
	Public  PublicKey
}

// [KeyRing] hold the set of local keys to use during handshake and session.
type KeyRing struct {
	kp DHKey     // encrypt-decrypt key pair exchange
	sv EDKeyPair // ED25519 local sign-verify keys
}

// [CipherState] provides symmetric encryption and decryption after a successful handshake.
// Please see [docs] for more information.
//
// [docs]: http://www.noiseprotocol.org/noise.html#the-cipherstate-object
type CipherState = *noise.CipherState

// [BytePool] implements a leaky pool of []byte in the form of a bounded channel.
type BytePool = *bpool.BytePool

// [HandshakeState] tracks the state of a Noise handshake.
// It may be discarded after the handshake is complete.
// Please see [docs] for more information.
//
// [docs]: http://www.noiseprotocol.org/noise.html#the-handshakestate-object
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
const headerSize = 2

// [CipherSuite] is a set of cryptographic primitives used in a Noise protocol.
// Based on: Diffie-Hellman X25519, [Blake2] and [ChaCha20-Poly1305]
// Please see [NoisePatternExplorer] for more details.
//
// [ChaCha20-Poly1305]: https://en.wikipedia.org/wiki/ChaCha20-Poly1305
// [Diffie-Hellman X25519]: https://en.wikipedia.org/wiki/Curve25519
// [Blake2]: https://www.blake2.net/

var CipherSuite = noise.NewCipherSuite(noise.DH25519, noise.CipherChaChaPoly, noise.HashBLAKE2b)

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
func newDHKeyPair() (DHKey, error) {
	// Diffie-Hellman key pair
	var err error
	var kp DHKey

	kp, err = noise.DH25519.GenerateKeypair(rand.Reader)
	if err != nil {
		err := fmt.Errorf("error trying to generate `s` keypair: %w", err)
		return kp, errDuringHandshake(err)
	}

	log.Printf("generated X25519 public key")
	return kp, nil
}

// newED25519KeyPair generate a new kp for sign-verify using P256 128 bit
func newED25519KeyPair() (EDKeyPair, error) {
	// ref: https://github.com/openssl/openssl/issues/18448
	// ref: https://csrc.nist.gov/csrc/media/events/workshop-on-elliptic-curve-cryptography-standards/documents/papers/session6-adalier-mehmet.pdf
	// TODO should i persist seed to keep the same public key for peer identity?
	// TODO rand.Reader store it and retrieve it to avoid change the pub key in every new handshake?
	pb, pv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return EDKeyPair{}, err
	}

	log.Print("generated ECDSA25519 public key")
	return EDKeyPair{pv, pb}, nil
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

// newHandshakeConfig create a noise config.
// A Config provides the details necessary to process a Noise handshake.
// It is never modified by this package, and can be reused.
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
	kr KeyRing
	hs HandshakeState
	p  BytePool
	i  bool
}

// newKeyRing create a bundle of local keys needed during session + handshake.
func newKeyRing() (KeyRing, error) {
	sv, err := newED25519KeyPair()
	if err != nil {
		return KeyRing{}, err
	}

	kp, err := newDHKeyPair()
	if err != nil {
		return KeyRing{}, err
	}

	return KeyRing{kp, sv}, nil
}

// newHandshake create a new handshake handler using provided connection and role.
func newHandshake(conn net.Conn, initiator bool) (*handshake, error) {
	kr, err := newKeyRing()
	if err != nil {
		return nil, err
	}

	// set handshake state as initiator?
	conf := newHandshakeConfig(initiator, kr.kp)
	// A HandshakeState tracks the state of a Noise handshake
	state, err := newHandshakeState(conf)
	if err != nil {
		return nil, errDuringHandshake(err)
	}

	// Setup the max of size possible for tokens exchanged between peers.
	edKeyLen := ed25519.PublicKeySize          // 32 bytes
	dhKeyLen := 2 * noise.DH25519.DHLen()      // 64 bytes
	cipherLen := 2 * chacha20poly1305.Overhead // 32 bytes
	// Sum the needed memory size for pool
	size := dhKeyLen + edKeyLen + cipherLen + headerSize
	pool := bpool.NewBytePool(bPools, size) // N bytes pool

	// Create a new session handler
	session, err := newSession(conn, kr)
	if err != nil {
		// Fail creating new session
		return nil, errDuringHandshake(err)
	}

	return &handshake{session, kr, state, pool, initiator}, nil
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

// Start initialize handshake based on peer rol
// If peer is initiator then Initiate method run else Answer
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
	log.Print("sending e to remote")
	enc, dec, err := h.Send()
	if err != nil {
		err = fmt.Errorf("error sending `e` state: %v", err)
		return errDuringHandshake(err)
	}

	// Receive message #2 stage
	log.Print("waiting for e, ee, s, es from remote")
	enc, dec, err = h.Receive()
	if err != nil {
		err = fmt.Errorf("error receiving `e, ee, s, es` state: %v", err)
		return errDuringHandshake(err)
	}

	// Send last handshake message #3 stage
	log.Print("sending s, se to remote")
	enc, dec, err = h.Send()
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
	// Receive message #1 stage
	log.Print("waiting for e from remote")
	enc, dec, err := h.Receive()
	if err != nil {
		err = fmt.Errorf("error receiving `e` state: %v", err)
		return errDuringHandshake(err)
	}

	// Send answer message #2 stage
	log.Print("sending e, ee, s, es to remote")
	enc, dec, err = h.Send()
	if err != nil {
		err = fmt.Errorf("error sending `e, ee, s, es` state: %v", err)
		return errDuringHandshake(err)
	}

	// Receive message #2 stage
	log.Print("waiting for s, se from remote")
	enc, dec, err = h.Receive()
	if err != nil {
		err = fmt.Errorf("error receiving `s, se` state: %v", err)
		return errDuringHandshake(err)
	}

	// Check if synchronization is valid
	if err := h.Valid(enc, dec); err != nil {
		return err
	}

	// Add keys for encrypt/decrypt operations in session.
	h.s.SetCyphers(dec, enc)
	return nil
}

// Send create a new token based on message pattern synchronization and send it to remote peer.
func (h *handshake) Send() (e, d CipherState, err error) {
	var msg []byte
	// Get a chunk of bytes from pool
	// We need an empty buffer slice here
	buffer := h.p.Get()[:0]
	defer h.p.Put(buffer)

	// WriteMessage appends a handshake message to out. The message will include the
	// optional payload if provided. If the handshake is completed by the call, two
	// CipherStates will be returned, one is used for encryption of messages to the
	// remote peer, the other is used for decryption of messages from the remote
	// peer. Append public signature key in payload to share with remote.
	msg, e, d, err = h.hs.WriteMessage(buffer, h.kr.sv.Public)
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
	var payload []byte // sent payload
	var size uint16    // read bytes size from header

	// Read incoming message size
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
	payload, e, d, err = h.hs.ReadMessage(nil, buffer)
	// Set remote signature validation public key
	h.s.SetRemotePublicKey(payload)
	return

}
