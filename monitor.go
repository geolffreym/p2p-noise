package noise

// 1- keep some list of peers and some measure of their "usability" - health, ping response, time connected, reliability, latency = more distance
// 2- prioritize staying connected to the top 10 (arbitrary, could be 5 could be 20
// 3- Connect to another set of peers more casually

// PC = ping request counts
// PR = ping response count
// RE = reliability
// LA = distance based on time of response

// health = PR / PC
// RE = last time ping - first time ping

// which is a good  peer?
// 1- the closest to my node
// 2- healthy node answering
// 3- long time peer connected

// Discovery

// Connect to peer
// Peer send me his healthy node A calculated
// my peer connect with A node and receive his healthy node B
// my peer connect with B node and receive his healthy node C
// .... son until i reach the limit of connected nodes
