package interfaces

import (
	"os"
	"testing"
)

func requireIntegration(t *testing.T) {
	t.Helper()
	if os.Getenv("OPNSENSE_URI") == "" || os.Getenv("OPNSENSE_API_KEY") == "" || os.Getenv("OPNSENSE_API_SECRET") == "" {
		t.Skip("OPNSENSE_URI/OPNSENSE_API_KEY/OPNSENSE_API_SECRET not set; skipping integration test")
	}
}
