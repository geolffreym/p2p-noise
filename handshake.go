package noise

// https://en.wikipedia.org/wiki/Handshaking
// http://www.noiseprotocol.org/noise.html
// In telecommunications, a handshake is an automated process of
// negotiation between two participants (example "Alice and Bob")
// through the exchange of information that establishes the protocols of a
// communication link at the start of the communication, before full communication begins

type Noise struct {
}

func NewHandshake() {

}

// Add Chain of Responsibility for handshake events
// Allow plugable middlewares
