// Package aggregator provides a thread-safe mechanism to combine port scan
// results from multiple named sources into a single deduplicated port list.
//
// Each source is identified by a string key and can be updated independently.
// Calling Merge returns a unified view of all currently known open ports,
// with duplicates (same port number and protocol) removed.
//
// Typical usage:
//
//	agg := aggregator.New()
//	agg.Update("tcp", tcpPorts)
//	agg.Update("udp", udpPorts)
//	allPorts := agg.Merge()
package aggregator
