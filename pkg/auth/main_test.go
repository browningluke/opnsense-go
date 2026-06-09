package auth

import (
	"fmt"
	"os"
	"testing"
)

// Integration tests in this package require a live OPNsense instance. Skip
// cleanly when the required env vars are unset so `go test ./...` doesn't
// fail with cryptic network errors on a developer's first checkout.
func TestMain(m *testing.M) {
	if os.Getenv("OPNSENSE_URI") == "" || os.Getenv("OPNSENSE_API_KEY") == "" || os.Getenv("OPNSENSE_API_SECRET") == "" {
		fmt.Fprintln(os.Stderr, "OPNSENSE_URI/OPNSENSE_API_KEY/OPNSENSE_API_SECRET not set; skipping integration tests in pkg/auth")
		os.Exit(0)
	}
	os.Exit(m.Run())
}
