package portpolicy

import (
	"fmt"
	"testing"

	"portwatch/internal/scanner"
	"portwatch/internal/state"
)

func BenchmarkEvaluate_10Policies_10Ports(b *testing.B) {
	e := New()
	for i := 0; i < 10; i++ {
		_ = e.Add(Policy{
			Name:     fmt.Sprintf("policy-%d", i),
			MinPort:  i * 1000,
			MaxPort:  i*1000 + 999,
			OnAdded:  true,
			Severity: SeverityWarn,
		})
	}

	ports := make([]scanner.Port, 10)
	for i := range ports {
		ports[i] = scanner.Port{Number: i * 1000, Protocol: "tcp"}
	}
	diff := state.Diff{Added: ports}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Evaluate(diff)
	}
}

func BenchmarkEvaluate_100Policies_NoPorts(b *testing.B) {
	e := New()
	for i := 0; i < 100; i++ {
		_ = e.Add(Policy{
			Name:     fmt.Sprintf("policy-%d", i),
			MinPort:  i * 100,
			MaxPort:  i*100 + 99,
			OnAdded:  true,
			Severity: SeverityWarn,
		})
	}
	diff := state.Diff{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Evaluate(diff)
	}
}
