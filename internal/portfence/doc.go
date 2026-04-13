// Package portfence enforces port-level access policies for portwatch.
//
// A Fence holds an allowlist and a blocklist of ports. When operating in
// strict mode every observed port must appear in the allowlist or a
// ViolationNotAllowed is raised. Regardless of mode, any port present in
// the blocklist always produces a ViolationBlocked.
//
// Typical usage:
//
//	f := portfence.New(true) // strict mode
//	f.Allow(scanner.Port{Number: 443, Protocol: "tcp"})
//	f.Block(scanner.Port{Number: 23,  Protocol: "tcp"})
//
//	violations := f.Check(observedPorts)
//	for _, v := range violations {
//		log.Println(v)
//	}
package portfence
