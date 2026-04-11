// Package eventlog implements a bounded, thread-safe append-only log of
// port change events for portwatch.
//
// Events are stored in insertion order and can be queried by time range,
// port identity, or event kind ("added" / "removed"). The log evicts the
// oldest entry when it reaches its configured maximum size, keeping memory
// usage predictable during long-running daemon sessions.
//
// Persistence is provided via Save and Load, which use newline-delimited
// JSON so the file remains human-readable and easy to tail.
package eventlog
