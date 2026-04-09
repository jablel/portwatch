package history

import (
	"encoding/json"
	"os"
	"time"

	"github.com/user/portwatch/internal/state"
)

// Entry represents a single historical snapshot with a timestamp.
type Entry struct {
	Timestamp time.Time    `json:"timestamp"`
	Ports     []state.Port `json:"ports"`
}

// History holds an ordered list of scan entries.
type History struct {
	Entries []Entry `json:"entries"`
	maxSize int
}

// New creates a History with the given maximum number of retained entries.
func New(maxSize int) *History {
	if maxSize <= 0 {
		maxSize = 100
	}
	return &History{maxSize: maxSize}
}

// Add appends a new entry, evicting the oldest if capacity is exceeded.
func (h *History) Add(ports []state.Port) {
	entry := Entry{Timestamp: time.Now().UTC(), Ports: ports}
	h.Entries = append(h.Entries, entry)
	if len(h.Entries) > h.maxSize {
		h.Entries = h.Entries[len(h.Entries)-h.maxSize:]
	}
}

// Save writes the history to a JSON file at the given path.
func (h *History) Save(path string) error {
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// Load reads history from a JSON file. Returns an empty History if the file
// does not exist.
func Load(path string, maxSize int) (*History, error) {
	h := New(maxSize)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return h, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, h); err != nil {
		return nil, err
	}
	return h, nil
}

// Latest returns the most recent entry, or nil if history is empty.
func (h *History) Latest() *Entry {
	if len(h.Entries) == 0 {
		return nil
	}
	e := h.Entries[len(h.Entries)-1]
	return &e
}
