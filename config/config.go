// Package conf provide a "functional option" design pattern to handle node settings.
// See also: https://github.com/crazybber/awesome-patterns/blob/master/idiom/functional-options.md
package config

import "time"

// Functional options
type Config struct {
	protocol             string
	selfListeningAddress string
	maxPayloadSize       uint32
	maxPeersConnected    uint8
	dialTimeout          time.Duration
	peerDeadline         time.Duration
}

type Setter func(*Config)

// Return default settings
func New() *Config {
	return &Config{
		// default protocol
		protocol: "tcp",
		// Self listening address
		selfListeningAddress: "0.0.0.0:8010",
		// Max payload size received from peers
		maxPayloadSize: 10 << 20, // 10MB
		// Max peer consecutively connected.
		// Each of this peers is equivalent to one routine, limit this is a performance consideration.
		maxPeersConnected: 100,
		// Max time waiting for dial to complete.
		// Default 5 seconds
		// ref: https://pkg.go.dev/net#DialTimeout
		dialTimeout: 5 * time.Second,
		// Max time waiting for I/O or peer interaction. After this time the connection will timeout and considered inactive.
		// Default 1800 seconds = 30 minutes.
		peerDeadline: 1800,
	}
}

// Write stores settings in `Settings` struct reference.
// All the settings are passed as an array of `setters` to then get called with `Settings“ reference as param.
// ref: https://github.com/crazybber/awesome-patterns/blob/master/idiom/functional-options.md
func (s *Config) Write(c ...Setter) {
	for _, setter := range c {
		setter(s)
	}
}

// Protocol returns the protocol to use for communication.
func (s *Config) Protocol() string {
	return s.protocol
}

// SelfListeningAddress returns the local node address.
func (s *Config) SelfListeningAddress() string {
	return s.selfListeningAddress
}

// MaxPeersConnected returns the max number of connections.
func (s *Config) MaxPeersConnected() uint8 {
	return s.maxPeersConnected
}

// MaxPayloadSize returns the max payload size allowed to received from peers.
func (s *Config) MaxPayloadSize() uint32 {
	return s.maxPayloadSize
}

// DialTimeOut returns max time waiting for dial to complete.
func (s *Config) DialTimeout() time.Duration {
	return s.dialTimeout
}

// PeerDeadline returns the max time waiting for I/O or peer interaction
func (s *Config) PeerDeadline() time.Duration {
	return s.peerDeadline
}

// Protocol sets the protocol to use when communicating.
func SetProtocol(protocol string) Setter {
	return func(conf *Config) {
		conf.protocol = protocol
	}
}

// SetSelfListeningAddress sets the listening address for local node.
func SetSelfListeningAddress(address string) Setter {
	return func(conf *Config) {
		conf.selfListeningAddress = address
	}
}

// SetMaxPeersConnected sets the maximum number of connections allowed for routing.
// If the number of connections > MaxPeersConnected then router drop new connections.
func SetMaxPeersConnected(maxPeers uint8) Setter {
	return func(conf *Config) {
		conf.maxPeersConnected = maxPeers
	}
}

// SetMaxPayloadSize sets the maximum bytes size received from peers.
// If the size exceed > MaxPayloadSize then payload is dropped.
func SetMaxPayloadSize(maxPayloadSize uint32) Setter {
	return func(conf *Config) {
		conf.maxPayloadSize = maxPayloadSize
	}
}

// SetPeerDeadline sets how long in seconds peer connections can remain idle.
// If deadline for I/O is exceeded a timeout is raised and the connection is closed.
// A deadline is an absolute time after which I/O operations
// fail instead of blocking. The deadline applies to all future
// and pending I/O, not just the immediately following call to Read or Write.
// ref: https://pkg.go.dev/net#Conn
func SetPeerDeadline(timeout time.Duration) Setter {
	return func(conf *Config) {
		conf.peerDeadline = timeout
	}
}

// SetDialTimeOut sets how long time in seconds the node will wait for dial timeouts.
func SetDialTimeout(timeout time.Duration) Setter {
	return func(conf *Config) {
		conf.dialTimeout = timeout
	}
}
