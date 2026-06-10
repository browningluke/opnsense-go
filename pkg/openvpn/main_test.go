package openvpn

import (
	"fmt"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func newTestController(t *testing.T) *Controller {
	t.Helper()

	client := api.NewClient(api.Options{
		Uri:           os.Getenv("OPNSENSE_URI"),
		APIKey:        os.Getenv("OPNSENSE_API_KEY"),
		APISecret:     os.Getenv("OPNSENSE_API_SECRET"),
		AllowInsecure: true,
		MaxBackoff:    30,
		MinBackoff:    1,
		MaxRetries:    4,
	})

	return &Controller{Api: client}
}

func TestMain(m *testing.M) {
	if os.Getenv("OPNSENSE_URI") == "" || os.Getenv("OPNSENSE_API_KEY") == "" || os.Getenv("OPNSENSE_API_SECRET") == "" {
		fmt.Fprintln(os.Stderr, "OPNSENSE_URI/OPNSENSE_API_KEY/OPNSENSE_API_SECRET not set; skipping integration tests in pkg/openvpn")
		os.Exit(0)
	}
	os.Exit(m.Run())
}
