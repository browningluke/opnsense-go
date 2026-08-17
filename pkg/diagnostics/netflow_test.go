package diagnostics

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestNetflow(t *testing.T) {
	apiClient := api.NewClient(api.Options{
		Uri:           os.Getenv("OPNSENSE_URI"),
		APIKey:        os.Getenv("OPNSENSE_API_KEY"),
		APISecret:     os.Getenv("OPNSENSE_API_SECRET"),
		AllowInsecure: true,
		MaxBackoff:    30,
		MinBackoff:    1,
		MaxRetries:    4,
	})

	controller := Controller{Api: apiClient}
	ctx := context.Background()

	// Step 1: Get current config (equivalent to import step).
	respGet, err := controller.NetflowGetConfig(ctx)
	if err != nil {
		t.Fatalf("NetflowGetConfig: %v", err)
	}
	t.Logf("NetflowGetConfig (initial): %+v", respGet)

	origActiveTimeout := respGet.Netflow.ActiveTimeout
	origInactiveTimeout := respGet.Netflow.InactiveTimeout

	// Step 2: Update timeouts (safe fields; doesn't toggle capture/collection).
	updated := respGet.Netflow
	updated.ActiveTimeout = "1801"
	updated.InactiveTimeout = "16"

	respSet, err := controller.NetflowSetConfig(ctx, &updated)
	if err != nil {
		t.Fatalf("NetflowSetConfig: %v", err)
	}
	if respSet.Result != "saved" {
		t.Fatalf("NetflowSetConfig: got result %q, want %q", respSet.Result, "saved")
	}

	if _, err := controller.NetflowReconfigure(ctx); err != nil {
		t.Fatalf("NetflowReconfigure: %v", err)
	}

	// Step 3: Read back and verify.
	respGet, err = controller.NetflowGetConfig(ctx)
	if err != nil {
		t.Fatalf("NetflowGetConfig (after update): %v", err)
	}
	if respGet.Netflow.ActiveTimeout != "1801" {
		t.Fatalf("activeTimeout not updated; got %q, want %q", respGet.Netflow.ActiveTimeout, "1801")
	}
	if respGet.Netflow.InactiveTimeout != "16" {
		t.Fatalf("inactiveTimeout not updated; got %q, want %q", respGet.Netflow.InactiveTimeout, "16")
	}

	// Step 4: Restore originals.
	restored := respGet.Netflow
	restored.ActiveTimeout = origActiveTimeout
	restored.InactiveTimeout = origInactiveTimeout
	if _, err := controller.NetflowSetConfig(ctx, &restored); err != nil {
		t.Fatalf("NetflowSetConfig (restore): %v", err)
	}
	if _, err := controller.NetflowReconfigure(ctx); err != nil {
		t.Fatalf("NetflowReconfigure (restore): %v", err)
	}

	// Step 5: Read-only RPCs.
	if _, err := controller.NetflowIsEnabled(ctx); err != nil {
		t.Fatalf("NetflowIsEnabled: %v", err)
	}
	if _, err := controller.NetflowStatus(ctx); err != nil {
		t.Fatalf("NetflowStatus: %v", err)
	}
	if _, err := controller.NetflowCacheStats(ctx); err != nil {
		t.Fatalf("NetflowCacheStats: %v", err)
	}
}
