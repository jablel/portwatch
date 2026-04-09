package reporter_test

import (
	"testing"

	"github.com/user/portwatch/internal/reporter"
)

func TestParseFormat_Valid(t *testing.T) {
	cases := []struct {
		input    string
		expected reporter.Format
	}{
		{"text", reporter.FormatText},
		{"csv", reporter.FormatCSV},
		{"", reporter.FormatText},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := reporter.ParseFormat(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := reporter.ParseFormat("json")
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}

func TestFormat_String(t *testing.T) {
	if reporter.FormatText.String() != "text" {
		t.Errorf("expected 'text', got %q", reporter.FormatText.String())
	}
	if reporter.FormatCSV.String() != "csv" {
		t.Errorf("expected 'csv', got %q", reporter.FormatCSV.String())
	}
}
