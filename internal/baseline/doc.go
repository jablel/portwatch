// Package baseline provides the trusted port snapshot used by portwatch
// to distinguish expected from unexpected open ports.
//
// Workflow:
//
//  1. On first run (or after an explicit "accept" command) the current
//     scan result is written to disk as the baseline.
//
//  2. On every subsequent scan the live port list is compared against
//     the baseline; any port not present in the baseline is treated as
//     an anomaly and forwarded to the alert/notifier pipeline.
//
//  3. The operator can update the baseline at any time by running
//     portwatch with the --update-baseline flag, which overwrites the
//     file with the current scan result.
//
// The baseline is stored as a JSON file (default: ~/.portwatch/baseline.json)
// and is safe to commit to version control for auditing purposes.
package baseline
