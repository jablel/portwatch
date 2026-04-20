package portdrift_test

import (
	"testing"

	"github.com/user/portwatch/internal/portdrift"
	"github.com/user/portwatch/internal/scanner"
)

func makePorts(specs [][2]interface{}) []scanner.Port {
	ports := make([]scanner.Port, 0, len(specs))
	for _, s := range specs {
		ports = append(ports, scanner.Port{
			Number:   s[0].(int),
			Protocol: s[1].(string),
		})
	}
	return ports
}

func TestScore_NoBaseline_ReturnsZero(t *testing.T) {
	tr := portdrift.New()
	current := makePorts([][2]interface{}{{80, "tcp"}, {443, "tcp"}})
	if got := tr.Score(current); got != 0.0 {
		t.Fatalf("expected 0.0 before baseline set, got %f", got)
	}
}

func TestScore_IdenticalSets_ReturnsZero(t *testing.T) {
	tr := portdrift.New()
	ports := makePorts([][2]interface{}{{80, "tcp"}, {443, "tcp"}})
	tr.SetBaseline(ports)
	if got := tr.Score(ports); got != 0.0 {
		t.Fatalf("expected 0.0 for identical sets, got %f", got)
	}
}

func TestScore_CompletelyDifferent_ReturnsOne(t *testing.T) {
	tr := portdrift.New()
	baseline := makePorts([][2]interface{}{{80, "tcp"}})
	current := makePorts([][2]interface{}{{9000, "tcp"}})
	tr.SetBaseline(baseline)
	if got := tr.Score(current); got != 1.0 {
		t.Fatalf("expected 1.0 for disjoint sets, got %f", got)
	}
}

func TestScore_HalfOverlap(t *testing.T) {
	tr := portdrift.New()
	baseline := makePorts([][2]interface{}{{80, "tcp"}, {443, "tcp"}})
	current := makePorts([][2]interface{}{{80, "tcp"}, {8080, "tcp"}})
	tr.SetBaseline(baseline)
	// union={80,443,8080}=3, intersect={80}=1 → drift = 1 - 1/3 ≈ 0.666
	got := tr.Score(current)
	if got < 0.66 || got > 0.67 {
		t.Fatalf("expected ~0.666, got %f", got)
	}
}

func TestScore_EmptyCurrentAgainstBaseline(t *testing.T) {
	tr := portdrift.New()
	baseline := makePorts([][2]interface{}{{80, "tcp"}})
	tr.SetBaseline(baseline)
	if got := tr.Score(nil); got != 1.0 {
		t.Fatalf("expected 1.0 when current is empty, got %f", got)
	}
}

func TestHasBaseline_FalseBeforeSet(t *testing.T) {
	tr := portdrift.New()
	if tr.HasBaseline() {
		t.Fatal("expected HasBaseline to be false before SetBaseline")
	}
}

func TestHasBaseline_TrueAfterSet(t *testing.T) {
	tr := portdrift.New()
	tr.SetBaseline(nil)
	if !tr.HasBaseline() {
		t.Fatal("expected HasBaseline to be true after SetBaseline")
	}
}

func TestSetBaseline_OverwritesPrevious(t *testing.T) {
	tr := portdrift.New()
	old := makePorts([][2]interface{}{{80, "tcp"}})
	newBase := makePorts([][2]interface{}{{9000, "tcp"}})
	tr.SetBaseline(old)
	tr.SetBaseline(newBase)
	// current matches new baseline → drift should be 0
	if got := tr.Score(newBase); got != 0.0 {
		t.Fatalf("expected 0.0 after overwriting baseline, got %f", got)
	}
}
