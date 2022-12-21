package config

import (
	"testing"
	"time"
)

func TestWrite(t *testing.T) {

	settings := []struct {
		Name              string
		MaxPeersConnected uint8
	}{{
		Name:              "10",
		MaxPeersConnected: 10,
	}, {
		Name:              "20",
		MaxPeersConnected: 20,
	}, {
		Name:              "30",
		MaxPeersConnected: 30,
	}, {
		Name:              "60",
		MaxPeersConnected: 60,
	}, {
		Name:              "255",
		MaxPeersConnected: 255,
	}}

	myLib := func(c ...Setter) *Config {
		s := New()
		s.Write(c...)
		return s
	}

	for _, e := range settings {
		t.Run(e.Name, func(t *testing.T) {
			libWithSettings := myLib(SetMaxPeersConnected(e.MaxPeersConnected))

			if libWithSettings.MaxPeersConnected() != e.MaxPeersConnected {
				t.Errorf("expected MaxPeerConnected %#v, get settings %v", e.MaxPeersConnected, libWithSettings.MaxPeersConnected())
			}

		})

	}

}

func TestSetMaxPeersConnected(t *testing.T) {
	settings := New()
	callable := SetMaxPeersConnected(10)
	callable(settings)

	if settings.MaxPeersConnected() != 10 {
		t.Errorf("expected MaxPeerConnected %#v, got settings %v", 10, settings.MaxPeersConnected())
	}
}

func TestPeerDeadline(t *testing.T) {
	settings := New()
	callable := SetIdleTimeout(100)
	callable(settings)

	if settings.IdleTimeout() != 100 {
		t.Errorf("expected MaxPeerConnected %#v, got settings %v", 10, settings.MaxPeersConnected())
	}
}

func TestMaxPayloadExceeded(t *testing.T) {
	settings := New()
	payloadSize := 1024
	callable := SetPoolBufferSize(payloadSize)
	callable(settings)

	if settings.PoolBufferSize() != payloadSize {
		t.Errorf("expected MaxPayloadExceeded %#v, got settings %v", payloadSize, settings.PoolBufferSize())
	}
}

func TestSelfListeningAddress(t *testing.T) {
	settings := New()
	address := "127.0.0.1:5003"
	callable := SetSelfListeningAddress(address)
	callable(settings)

	if settings.SelfListeningAddress() != address {
		t.Errorf("expected SelfListeningAddress %#v, got settings %v", address, settings.SelfListeningAddress())
	}
}

func TestProtocol(t *testing.T) {
	settings := New()
	protocol := "tcp"
	callable := SetProtocol(protocol)
	callable(settings)

	if settings.Protocol() != protocol {
		t.Errorf("expected Protocol %#v, got settings %v", protocol, settings.Protocol())
	}
}

func TestDialTimeout(t *testing.T) {
	settings := New()
	expected := time.Second * 10
	callable := SetDialTimeout(expected)
	callable(settings)

	if settings.DialTimeout() != expected {
		t.Errorf("expected DialTimeout %#v, got settings %v", expected, settings.DialTimeout())
	}
}
