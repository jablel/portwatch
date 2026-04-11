// Package tagger maps open ports to human-readable service labels.
//
// A Tagger consults a built-in table of well-known port-to-service mappings
// (e.g. 22→ssh, 80→http) and falls back to user-defined custom rules set
// via Define. When no mapping is found the label is formatted as
// "unknown:PORT".
//
// Usage:
//
//	tg := tagger.New()
//	tg.Define(8080, "my-api")        // optional custom override
//	label := tg.Tag(scanner.Port{Number: 443, Protocol: "tcp"}) // "https"
//	labels := tg.TagAll(ports)       // annotate a whole snapshot
//
// Tagger is safe for concurrent use.
package tagger
