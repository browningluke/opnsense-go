package crowdsec

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

// Integration tests in this package require a live OPNsense instance. Skip
// cleanly when the required env vars are unset so `go test ./...` doesn't
// fail with cryptic network errors on a developer's first checkout.
//
// CrowdSec also ships as the os-crowdsec plugin rather than in core, so its
// API is absent on instances that don't have it installed — including the CI
// image. Probe the endpoint and skip when it 404s, so an uninstalled plugin
// reads as "not applicable" rather than a test failure.
func TestMain(m *testing.M) {
	if os.Getenv("OPNSENSE_URI") == "" || os.Getenv("OPNSENSE_API_KEY") == "" || os.Getenv("OPNSENSE_API_SECRET") == "" {
		fmt.Fprintln(os.Stderr, "OPNSENSE_URI/OPNSENSE_API_KEY/OPNSENSE_API_SECRET not set; skipping integration tests in pkg/crowdsec")
		os.Exit(0)
	}

	if _, err := newController().GeneralGet(context.Background()); err != nil {
		if strings.Contains(err.Error(), "404") {
			fmt.Fprintln(os.Stderr, "os-crowdsec plugin not installed (404 from /crowdsec/general/get); skipping integration tests in pkg/crowdsec")
			os.Exit(0)
		}
	}

	os.Exit(m.Run())
}

func newController() *Controller {
	return &Controller{
		Api: api.NewClient(api.Options{
			Uri:           os.Getenv("OPNSENSE_URI"),
			APIKey:        os.Getenv("OPNSENSE_API_KEY"),
			APISecret:     os.Getenv("OPNSENSE_API_SECRET"),
			AllowInsecure: true,
			MaxBackoff:    30,
			MinBackoff:    1,
			MaxRetries:    4,
		}),
	}
}
