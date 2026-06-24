// Package install registers all in-source integrations for use.
//
// Integration registration is split by build profile:
//   - profile_full.go  (!alloy_slim) — full upstream-style integrations
//   - profile_slim.go  (alloy_slim)  — metrics-focused slim integrations
package install
