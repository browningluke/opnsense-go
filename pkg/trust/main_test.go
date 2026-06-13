package trust

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if os.Getenv("OPNSENSE_URI") == "" || os.Getenv("OPNSENSE_API_KEY") == "" || os.Getenv("OPNSENSE_API_SECRET") == "" {
		fmt.Fprintln(os.Stderr, "OPNSENSE_URI/OPNSENSE_API_KEY/OPNSENSE_API_SECRET not set; skipping integration tests in pkg/trust")
		os.Exit(0)
	}
	os.Exit(m.Run())
}
