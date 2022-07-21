package config

import (
	"testing"
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

	myLib := func(c ...Config) *Configs {
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
		t.Errorf("expected MaxPeerConnected %#v, get settings %v", 10, settings.MaxPeersConnected())
	}
}

func TestPeerDeadline(t *testing.T) {
	settings := New()
	callable := SetPeerDeadline(100)
	callable(settings)

	if settings.PeerDeadline() != 100 {
		t.Errorf("expected MaxPeerConnected %#v, get settings %v", 10, settings.MaxPeersConnected())
	}
}
