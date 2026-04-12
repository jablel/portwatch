// Package windowcounter implements a thread-safe sliding-window event
// counter keyed by an arbitrary string (typically a "protocol:port"
// identifier).
//
// Use it to track how frequently a particular port event fires over a
// rolling duration — for example, to detect flapping ports that open and
// close repeatedly within a short period.
//
// Example:
//
//	c := windowcounter.New(30 * time.Second)
//	count := c.Add("tcp:8080")  // returns total hits in the last 30 s
//	if count >= 5 {
//		// port 8080 has flapped at least 5 times in 30 seconds
//	}
package windowcounter
