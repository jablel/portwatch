// Package pipeline provides a single-step abstraction that composes the
// scanner, filter, state-diff, and notifier components into one reusable
// Run call.
//
// Typical usage:
//
//	sc  := scanner.New(host, ports, proto)
//	f   := filter.New(rules)
//	n, _ := notifier.New(notifier.BackendStdout, nil)
//	p   := pipeline.New(sc, f, n)
//
//	var previous []scanner.Port
//	for {
//		res, err := p.Run(ctx, previous)
//		if err != nil { log.Println(err) }
//		previous = res.Ports
//		time.Sleep(interval)
//	}
package pipeline
