package unbound

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestSettings(t *testing.T) {
	opnsense_url := os.Getenv("OPNSENSE_URI")
	opnsense_key := os.Getenv("OPNSENSE_API_KEY")
	opnsense_secret := os.Getenv("OPNSENSE_API_SECRET")

	api_client := api.NewClient(api.Options{
		Uri:           opnsense_url,
		APIKey:        opnsense_key,
		APISecret:     opnsense_secret,
		AllowInsecure: true,
		MaxBackoff:    30,
		MinBackoff:    1,
		MaxRetries:    4,
	})

	controller := Controller{
		Api: api_client,
	}
	ctx := context.Background()

	// Step 1: Get current settings (equivalent to import step).
	respGet, err := controller.SettingsGet(ctx)
	if err != nil {
		t.Fatalf("Failed to get settings: %v", err)
	}
	t.Logf("SettingsGet (initial): %+v", respGet)

	// Verify key fields are present in the response.
	if respGet.Unbound.General.Port == "" {
		t.Fatal("Expected general.port to be set")
	}
	if respGet.Unbound.General.Enabled == "" {
		t.Fatal("Expected general.enabled to be set")
	}
	if respGet.Unbound.ACLs.DefaultAction.String() == "" {
		t.Fatal("Expected acls.default_action to be set")
	}

	// Preserve originals so we can restore after the test.
	origHideIdentity := respGet.Unbound.Advanced.HideIdentity
	origHideVersion := respGet.Unbound.Advanced.HideVersion
	origLogQueries := respGet.Unbound.Advanced.LogQueries

	// Step 2: Update safe advanced settings (hide_identity, hide_version, log_queries -> true).
	updated := respGet.Unbound
	updated.Advanced.HideIdentity = "1"
	updated.Advanced.HideVersion = "1"
	updated.Advanced.LogQueries = "1"

	respSet, err := controller.SettingsUpdate(ctx, &updated)
	if err != nil {
		t.Fatalf("Failed to update settings: %v", err)
	}
	t.Logf("SettingsUpdate (enable flags): %+v", respSet)

	// Step 3: Read back and verify the changes were applied.
	respGet, err = controller.SettingsGet(ctx)
	if err != nil {
		t.Fatalf("Failed to get settings after update: %v", err)
	}
	t.Logf("SettingsGet (after update): %+v", respGet)

	if respGet.Unbound.Advanced.HideIdentity != "1" {
		t.Fatalf("hide_identity not updated; got %q, want %q", respGet.Unbound.Advanced.HideIdentity, "1")
	}
	if respGet.Unbound.Advanced.HideVersion != "1" {
		t.Fatalf("hide_version not updated; got %q, want %q", respGet.Unbound.Advanced.HideVersion, "1")
	}
	if respGet.Unbound.Advanced.LogQueries != "1" {
		t.Fatalf("log_queries not updated; got %q, want %q", respGet.Unbound.Advanced.LogQueries, "1")
	}

	// Step 4: Restore original values.
	restored := respGet.Unbound
	restored.Advanced.HideIdentity = origHideIdentity
	restored.Advanced.HideVersion = origHideVersion
	restored.Advanced.LogQueries = origLogQueries

	respSet, err = controller.SettingsUpdate(ctx, &restored)
	if err != nil {
		t.Fatalf("Failed to restore settings: %v", err)
	}
	t.Logf("SettingsUpdate (restore): %+v", respSet)

	// Step 5: Verify restore.
	respGet, err = controller.SettingsGet(ctx)
	if err != nil {
		t.Fatalf("Failed to get settings after restore: %v", err)
	}
	t.Logf("SettingsGet (after restore): %+v", respGet)

	if respGet.Unbound.Advanced.HideIdentity != origHideIdentity {
		t.Fatalf("hide_identity not restored; got %q, want %q", respGet.Unbound.Advanced.HideIdentity, origHideIdentity)
	}
	if respGet.Unbound.Advanced.HideVersion != origHideVersion {
		t.Fatalf("hide_version not restored; got %q, want %q", respGet.Unbound.Advanced.HideVersion, origHideVersion)
	}
	if respGet.Unbound.Advanced.LogQueries != origLogQueries {
		t.Fatalf("log_queries not restored; got %q, want %q", respGet.Unbound.Advanced.LogQueries, origLogQueries)
	}
}
