package retry

import (
	"context"
	"testing"
	"time"
)

func BenchmarkDo_ImmediateSuccess(b *testing.B) {
	r := fastRetryer(Default())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = r.Do(context.Background(), func() error { return nil })
	}
}

func BenchmarkDo_AlwaysFails(b *testing.B) {
	cfg := Config{MaxAttempts: 3, BaseDelay: time.Millisecond, MaxDelay: time.Millisecond}
	r := fastRetryer(cfg)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = r.Do(context.Background(), func() error { return errTemp })
	}
}
