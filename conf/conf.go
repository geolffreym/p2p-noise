package conf

import "time"

// Functional options
type Settings struct {
	MaxPeersConnected uint8
	PeerDeadline      time.Duration
}

type Setting func(*Settings)

// Return default settings
func NewSettings() *Settings {
	return &Settings{
		MaxPeersConnected: 100,
	}
}

// Write stores settings in `Settings` struct reference.
// All the settings are passed as an array of `setters` to then get called with `Settings`` reference as param.
// ref: https://github.com/crazybber/awesome-patterns/blob/master/idiom/functional-options.md
func (s *Settings) Write(c ...Setting) {
	for _, setter := range c {
		setter(s)
	}
}

// SetMaxPeersConnected sets the maximum number of connections allowed for routing.
// If the number of connections > MaxPeersConnected then router drop new connections.
func SetMaxPeersConnected(maxPeers uint8) Setting {
	return func(conf *Settings) {
		conf.MaxPeersConnected = maxPeers
	}
}
