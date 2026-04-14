// Package portindex provides a thread-safe, multi-key in-memory index
// over a collection of scanned ports.
//
// After each scan cycle the caller rebuilds the index via Build; subsequent
// reads via ByNumber, ByProtocol, and ByTag are served from the in-memory
// maps without any I/O.
//
// Typical usage:
//
//	idx := portindex.New()
//	idx.Build(latestPorts)
//	httpPorts := idx.ByTag("http")
//	tcpPorts  := idx.ByProtocol("tcp")
package portindex
