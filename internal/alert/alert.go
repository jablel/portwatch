package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/state"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Alert represents a single port change notification.
type Alert struct {
	Timestamp time.Time
	Level     Level
	Message   string
	Port      state.Port
}

// Notifier sends alerts to a configured output.
type Notifier struct {
	out io.Writer
}

// New creates a Notifier that writes to the given writer.
// If w is nil, os.Stdout is used.
func New(w io.Writer) *Notifier {
	if w == nil {
		w = os.Stdout
	}
	return &Notifier{out: w}
}

// Notify formats and writes alerts for the given diff.
func (n *Notifier) Notify(diff state.Diff) {
	for _, p := range diff.Added {
		n.write(Alert{
			Timestamp: time.Now(),
			Level:     LevelAlert,
			Message:   "port opened",
			Port:      p,
		})
	}
	for _, p := range diff.Removed {
		n.write(Alert{
			Timestamp: time.Now(),
			Level:     LevelWarn,
			Message:   "port closed",
			Port:      p,
		})
	}
}

func (n *Notifier) write(a Alert) {
	fmt.Fprintf(n.out, "[%s] %s — %s\n",
		a.Level,
		a.Timestamp.Format(time.RFC3339),
		fmt.Sprintf("%s %s", a.Message, a.Port),
	)
}
