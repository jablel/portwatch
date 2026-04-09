package reporter

import "fmt"

// ParseFormat converts a string to a Format, returning an error for unknown values.
func ParseFormat(s string) (Format, error) {
	switch Format(s) {
	case FormatText, FormatCSV:
		return Format(s), nil
	case "":
		return FormatText, nil
	default:
		return "", fmt.Errorf("unknown report format %q: must be one of [text csv]", s)
	}
}

// String returns the string representation of a Format.
func (f Format) String() string {
	return string(f)
}
