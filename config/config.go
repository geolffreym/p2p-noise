// Package conf provide a "functional option" design pattern to handle node settings.
// See also: https://github.com/crazybber/awesome-patterns/blob/master/idiom/functional-options.md
package config

import "time"

// Functional options
type Config struct {
	maxPeersConnected    uint8
	lingerTime           int
	poolBufferSize       int
	protocol             string
	selfListeningAddress string
	keepAlivePeriod      time.Duration
	dialTimeout          time.Duration
	idleTimeout          time.Duration
}

type Setter func(*Config)

// Return default settings
func New() *Config {
	return &Config{
		// default protocol
		protocol: "tcp",
		// Keep alive message time interval
		keepAlivePeriod: 1800 * time.Second,
		// Self listening address
		selfListeningAddress: "0.0.0.0:",
		// Max buffer pool size to handle incoming messages.
		poolBufferSize: 10 << 20, // 10MB
		// Max peer consecutively connected.
		// Each of this peers is equivalent to one routine, limit this is a performance consideration.
		maxPeersConnected: 100,
		// Max time waiting for dial to complete.
		// Default 5 seconds
		// ref: https://pkg.go.dev/net#DialTimeout
		dialTimeout: 5 * time.Second,
		// Max time waiting for I/O or peer interaction. After this time the connection will timeout and considered inactive.
		// After every received/send message a new deadline is refreshed using this value.
		// When the Keep Alive Interval is greater than the Idle Timeout, the BIG-IP system never sends TCP Keep-Alive packets as the connections are removed when reaching the TCP Idle Timeout .
		// Default 0 seconds = no deadline.
		idleTimeout: 0,
		// Discard unsent data after N seconds.
		// 0 means immediately discard unsent data after close.
		lingerTime: 0,
	}
}

// Write stores settings in `Settings` struct reference.
// All the settings are passed as an array of `setters` to then get called with `Settings“ reference as param.
// ref: https://github.com/crazybber/awesome-patterns/blob/master/idiom/functional-options.md
func (c *Config) Write(s ...Setter) {
	for _, setter := range s {
		setter(c)
	}
}

// Protocol returns the protocol to use for communication.
func (c *Config) Protocol() string {
	return c.protocol
}

// KeepAlive tells to node if should keep alive TCP connection.
func (c *Config) KeepAlive() time.Duration {
	return c.keepAlivePeriod
}

// Linger returns the time to wait to discard messages after close node.
func (c *Config) Linger() int {
	return c.lingerTime
}

// SelfListeningAddress returns the local node address.
func (c *Config) SelfListeningAddress() string {
	return c.selfListeningAddress
}

// MaxPeersConnected returns the max number of connections.
func (c *Config) MaxPeersConnected() uint8 {
	return c.maxPeersConnected
}

// PoolBufferSize returns the max payload size allowed to received from peers.
func (c *Config) PoolBufferSize() int {
	return c.poolBufferSize
}

// DialTimeOut returns max time waiting for dial to complete.
func (c *Config) DialTimeout() time.Duration {
	return c.dialTimeout
}

// IdleTimeout returns the max time waiting for I/O or peer interaction.
func (c *Config) IdleTimeout() time.Duration {
	return c.idleTimeout
}

// SetKeepAlive set the flag to keep alive or not the TCP connection.
func SetKeepAlive(ka time.Duration) Setter {
	return func(conf *Config) {
		conf.keepAlivePeriod = ka
	}
}

// SetProtocol sets the protocol to use when communicating.
func SetProtocol(protocol string) Setter {
	return func(conf *Config) {
		conf.protocol = protocol
	}
}

// SetLinger set the linger time in seconds to wait to discard messages after close node.
func SetLinger(linger int) Setter {
	return func(conf *Config) {
		conf.lingerTime = linger
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

// SetPoolBufferSize sets the maximum bytes size received from peers.
// If the size exceed > MaxPayloadSize then payload is dropped.
func SetPoolBufferSize(maxPayloadSize int) Setter {
	return func(conf *Config) {
		conf.poolBufferSize = maxPayloadSize
	}
}

// SetIdleTimeout sets how long in seconds peer connections can remain idle.
// If deadline for I/O is exceeded a timeout is raised and the connection is closed.
// A deadline is an absolute time after which I/O operations
// fail instead of blocking. The deadline applies to all future
// and pending I/O, not just the immediately following call to Read or Write.
// ref: https://pkg.go.dev/net#Conn
func SetIdleTimeout(timeout time.Duration) Setter {
	return func(conf *Config) {
		conf.idleTimeout = timeout
	}
}

// SetDialTimeOut sets how long time in seconds the node will wait for dial timeouts.
func SetDialTimeout(timeout time.Duration) Setter {
	return func(conf *Config) {
		conf.dialTimeout = timeout
	}
}
