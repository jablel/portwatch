// Package portaudit provides a lightweight, append-only audit log for
// port-change events detected by portwatch.
//
// Each [Entry] captures the wall-clock time, the [Kind] of change
// (added or removed), the affected [scanner.Port], and an optional
// actor string that identifies which component or user triggered the
// record.
//
// The [Log] is safe for concurrent use. Entries are kept in memory and
// can be persisted to a JSON file via [Log.Save] and restored via
// [Log.Load]. A non-zero maxSize passed to [New] bounds memory use by
// evicting the oldest entry whenever the cap is exceeded.
package portaudit
