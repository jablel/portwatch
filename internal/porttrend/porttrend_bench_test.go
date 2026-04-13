package porttrend

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func BenchmarkRecord_100Ports(b *testing.B) {
	tr := New(time.Minute)
	ports := make([]scanner.Port, 100)
	for i := range ports {
		ports[i] = scanner.Port{Protocol: "tcp", Number: i + 1}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Record(ports)
	}
}

func BenchmarkTrend(b *testing.B) {
	tr := New(time.Minute)
	p := scanner.Port{Protocol: "tcp", Number: 443}
	for i := 0; i < 200; i++ {
		tr.Record([]scanner.Port{p})
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tr.Trend(p)
	}
}
