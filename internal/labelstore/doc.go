// Package labelstore provides a thread-safe, file-backed key/value store
// that maps scanner.Port values to human-readable string labels.
//
// Labels are stored as a JSON object keyed by "<proto>:<port>" strings,
// making the file easy to inspect and edit by hand. The store is safe
// for concurrent use from multiple goroutines.
//
// Typical usage:
//
//	ls := labelstore.New("/var/lib/portwatch/labels.json")
//	_ = ls.Load()
//
//	p := scanner.Port{Number: 8080, Proto: "tcp"}
//	ls.Set(p, "dev-proxy")
//	label, ok := ls.Get(p)
//	_ = ls.Save()
package labelstore
