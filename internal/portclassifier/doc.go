// Package portclassifier categorises ports into well-known tiers.
//
// Ports are divided into three tiers:
//   - System (0–1023): privileged ports reserved by IANA.
//   - Registered (1024–49151): assigned to specific services.
//   - Dynamic (49152–65535): ephemeral / private use.
//
// A Classifier can also attach a human-readable label to each port
// based on a built-in table of common services.
package portclassifier
