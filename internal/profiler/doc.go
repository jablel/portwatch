// Package profiler records per-port scan latency samples and exposes
// aggregated statistics (count, total, min, max, mean) for each port key.
//
// A port key is typically a string of the form "<protocol>:<port>", e.g.
// "tcp:443". The Profiler is safe for concurrent use.
//
// Usage:
//
//	p := profiler.New()
//	start := time.Now()
//	// ... perform scan ...
//	p.Record(profiler.Sample{
//		Key:        "tcp:80",
//		Duration:   time.Since(start),
//		RecordedAt: time.Now(),
//	})
//	st, _ := p.Get("tcp:80")
//	fmt.Println(st.Mean())
package profiler
