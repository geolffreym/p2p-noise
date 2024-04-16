package noise

// metrics hold the statistics related to remote peers.
// We can add any method related to adaptive lookup logic here.
// Please see [docs] for more information.
//
// !IMPORTANT The order of byte size needed for each type in structs matter and impact the struct size.
// The fields are distributed in a way that ensures their alignment in 8-byte blocks.
// For instance on a 64-bit CPU, alignment blocks are 8 bytes.
//
//	0: [latency, bandwidth, nonce, sent], // 8 bytes
//	1: [recv, handshakeTime, 2bytespadding], // 8 bytes
//	2: [bytesRecv], // 8 bytes
//	3: [bytesSent], // 8 bytes
//
// [docs]: https://arxiv.org/pdf/1509.04417.pdf
type metrics struct {
	latency       uint16 // rtt in ms: 2bytes
	bandwidth     uint16 // remote peer bandwidth: 2bytes
	nonce         uint16 // nonce ordering factor: 2bytes
	sent          uint16 // counter sent messages: 2bytes
	recv          uint16 // counter received messages: 2bytes
	handshakeTime uint32 // how long took the handshake to complete.: 4bytes
	bytesRecv     uint64 // bytes received: 8bytes
	bytesSent     uint64 // bytes sent: 8 bytes
}

// TODO https://community.f5.com/t5/technical-articles/introducing-tcp-analytics/ta-p/290873
// calculate weight
// builder pattern?
