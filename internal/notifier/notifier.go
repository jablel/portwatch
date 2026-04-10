// Package notifier provides pluggable notification backends for portwatch.
package notifier

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/portwatch/internal/state"
)

// Backend represents a notification delivery method.
type Backend string

const (
	BackendStdout Backend = "stdout"
	BackendFile   Backend = "file"
)

// ParseBackend parses a string into a Backend, returning an error if unknown.
func ParseBackend(s string) (Backend, error) {
	switch strings.ToLower(s) {
	case string(BackendStdout):
		return BackendStdout, nil
	case string(BackendFile):
		return BackendFile, nil
	default:
		return "", fmt.Errorf("unknown notifier backend: %q", s)
	}
}

// Notifier sends notifications about port changes.
type Notifier struct {
	backend Backend
	w       io.Writer
}

// New creates a Notifier for the given backend. For BackendFile, path must be
// non-empty; for BackendStdout, path is ignored.
func New(backend Backend, path string) (*Notifier, error) {
	var w io.Writer
	switch backend {
	case BackendStdout:
		w = os.Stdout
	case BackendFile:
		if path == "" {
			return nil, fmt.Errorf("notifier: file backend requires a non-empty path")
		}
		f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			return nil, fmt.Errorf("notifier: open file: %w", err)
		}
		w = f
	default:
		return nil, fmt.Errorf("notifier: unsupported backend: %q", backend)
	}
	return &Notifier{backend: backend, w: w}, nil
}

// Notify writes a human-readable summary of the diff to the configured backend.
func (n *Notifier) Notify(diff state.Diff) error {
	if len(diff.Added) == 0 && len(diff.Removed) == 0 {
		return nil
	}
	for _, p := range diff.Added {
		if _, err := fmt.Fprintf(n.w, "[portwatch] PORT OPENED  %s\n", p); err != nil {
			return err
		}
	}
	for _, p := range diff.Removed {
		if _, err := fmt.Fprintf(n.w, "[portwatch] PORT CLOSED  %s\n", p); err != nil {
			return err
		}
	}
	return nil
}
