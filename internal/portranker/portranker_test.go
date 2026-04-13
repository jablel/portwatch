package portranker_test

import (
	"testing"

	"github.com/example/portwatch/internal/portranker"
	"github.com/example/portwatch/internal/scanner"
)

func makePort(number int, proto string) scanner.Port {
	return scanner.Port{Number: number, Protocol: proto}
}

func TestRank_SystemPortScoresHighest(t *testing.T) {
	r := portranker.New(portranker.DefaultWeights())
	ports := []scanner.Port{
		makePort(55000, "tcp"), // dynamic
		makePort(80, "tcp"),    // system
		makePort(8080, "tcp"),  // registered
	}

	entries := r.Rank(ports)

	if entries[0].Port.Number != 80 {
		t.Fatalf("expected port 80 first, got %d", entries[0].Port.Number)
	}
}

func TestRank_DynamicPortScoresLowest(t *testing.T) {
	r := portranker.New(portranker.DefaultWeights())
	ports := []scanner.Port{
		makePort(49200, "tcp"),
		makePort(443, "tcp"),
	}

	entries := r.Rank(ports)

	if entries[len(entries)-1].Port.Number != 49200 {
		t.Fatalf("expected dynamic port last, got %d", entries[len(entries)-1].Port.Number)
	}
}

func TestRank_TCPBonusApplied(t *testing.T) {
	w := portranker.Weights{ClassBonus: 0, DynamicPenalty: 0, TCPBonus: 3.0}
	r := portranker.New(w)
	ports := []scanner.Port{
		makePort(514, "udp"),
		makePort(514, "tcp"),
	}

	entries := r.Rank(ports)

	if entries[0].Port.Protocol != "tcp" {
		t.Fatalf("expected tcp first, got %s", entries[0].Port.Protocol)
	}
}

func TestRank_EmptyInput(t *testing.T) {
	r := portranker.New(portranker.DefaultWeights())
	entries := r.Rank(nil)
	if len(entries) != 0 {
		t.Fatalf("expected empty result, got %d entries", len(entries))
	}
}

func TestSetWeights_UpdatesScoring(t *testing.T) {
	r := portranker.New(portranker.DefaultWeights())
	r.SetWeights(portranker.Weights{ClassBonus: 0, DynamicPenalty: 0, TCPBonus: 0})

	ports := []scanner.Port{
		makePort(80, "tcp"),
		makePort(9999, "udp"),
	}

	entries := r.Rank(ports)
	// With zero weights, order is by port number ascending.
	if entries[0].Port.Number != 80 {
		t.Fatalf("expected port 80 first, got %d", entries[0].Port.Number)
	}
}

func TestRank_TiesBrokenByPortNumber(t *testing.T) {
	w := portranker.Weights{ClassBonus: 0, DynamicPenalty: 0, TCPBonus: 0}
	r := portranker.New(w)
	ports := []scanner.Port{
		makePort(300, "tcp"),
		makePort(100, "tcp"),
		makePort(200, "tcp"),
	}

	entries := r.Rank(ports)

	for i, want := range []int{100, 200, 300} {
		if entries[i].Port.Number != want {
			t.Fatalf("index %d: expected %d, got %d", i, want, entries[i].Port.Number)
		}
	}
}
