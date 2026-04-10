// Package snapshot provides named, persistent captures of open port sets.
//
// A snapshot is a point-in-time record of which ports were observed open,
// stored as a JSON file under a configurable directory. Snapshots can be
// used to compare the current port state against a previously known-good
// baseline, or to track changes across arbitrary named checkpoints.
//
// Usage:
//
//	store, err := snapshot.New("/var/lib/portwatch/snapshots")
//	if err != nil { ... }
//
//	// capture current state
//	if err := store.Save("pre-deploy", currentPorts); err != nil { ... }
//
//	// later, load and compare
//	snap, err := store.Load("pre-deploy")
//	if err != nil { ... }
//	diff := state.Compare(snap.Ports, currentPorts)
package snapshot
