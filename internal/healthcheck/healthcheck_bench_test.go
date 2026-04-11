package healthcheck_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/healthcheck"
)

func BenchmarkRecordScan(b *testing.B) {
	c := healthcheck.New(5 * time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.RecordScan()
	}
}

func BenchmarkStatus(b *testing.B) {
	c := healthcheck.New(5 * time.Second)
	c.RecordScan()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.Status()
	}
}

func BenchmarkStatusString(b *testing.B) {
	c := healthcheck.New(5 * time.Second)
	c.RecordScan()
	s := c.Status()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.String()
	}
}
