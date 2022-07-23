//Copyright (c) 2022, Geolffrey Mena <gmjun2000@gmail.com>
//
//P2P Noise Secure handshake.
//See also: http://www.noiseprotocol.org/noise.html#introduction
package noise

import "time"

type Config interface {
	MaxPeersConnected() uint8
	PeerDeadline() time.Duration
}

type Node struct {
	net   *Network
	noise *noise
}
