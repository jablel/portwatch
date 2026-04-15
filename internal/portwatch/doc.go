// Package portwatch provides a high-level Watcher that coordinates a single
// scan-diff-notify cycle used by the portwatch daemon.
//
// Usage:
//
//	sc := scanner.New(cfg)
//	n, _ := notifier.New(notifier.BackendStdout, os.Stdout)
//	w := portwatch.New(sc, n)
//
//	for {
//		res, err := w.Tick(ctx)
//		if err != nil {
//			log.Println(err)
//		}
//		time.Sleep(cfg.Interval)
//	}
//
// Watcher keeps track of the previously observed port set internally, so
// callers do not need to manage state between ticks.
package portwatch
