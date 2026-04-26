package kea

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestSettingsDHCPv4(t *testing.T) {
	api_client := api.NewClient(api.Options{
		Uri:           os.Getenv("OPNSENSE_URI"),
		APIKey:        os.Getenv("OPNSENSE_API_KEY"),
		APISecret:     os.Getenv("OPNSENSE_API_SECRET"),
		AllowInsecure: true,
		MaxBackoff:    30,
		MinBackoff:    1,
		MaxRetries:    4,
	})

	controller := Controller{Api: api_client}
	ctx := context.Background()

	// Step 1: Get current settings.
	respGet, err := controller.SettingsDHCPv4Get(ctx)
	if err != nil {
		t.Fatalf("SettingsDHCPv4Get failed: %v", err)
	}
	t.Logf("SettingsDHCPv4Get (initial): %+v", respGet)

	// Verify key fields are present.
	if respGet.DHCPv4.General.ValidLifetime == "" {
		t.Fatal("Expected general.valid_lifetime to be set")
	}
	if respGet.DHCPv4.General.Enabled == "" {
		t.Fatal("Expected general.enabled to be set")
	}

	// Preserve original so we can restore after the test.
	origValidLifetime := respGet.DHCPv4.General.ValidLifetime
	origFWRules := respGet.DHCPv4.General.FWRules

	// Step 2: Update safe fields.
	updated := respGet.DHCPv4
	updated.General.ValidLifetime = "3600"
	updated.General.FWRules = "0"

	respSet, err := controller.SettingsDHCPv4Update(ctx, &updated)
	if err != nil {
		t.Fatalf("SettingsDHCPv4Update failed: %v", err)
	}
	t.Logf("SettingsDHCPv4Update (update): %+v", respSet)
	if respSet.Result != "saved" {
		t.Fatalf("Expected result 'saved', got %q", respSet.Result)
	}

	// Step 3: Read back and verify.
	respGet, err = controller.SettingsDHCPv4Get(ctx)
	if err != nil {
		t.Fatalf("SettingsDHCPv4Get (after update) failed: %v", err)
	}
	if respGet.DHCPv4.General.ValidLifetime != "3600" {
		t.Fatalf("valid_lifetime not updated; got %q, want %q", respGet.DHCPv4.General.ValidLifetime, "3600")
	}
	if respGet.DHCPv4.General.FWRules != "0" {
		t.Fatalf("fwrules not updated; got %q, want %q", respGet.DHCPv4.General.FWRules, "0")
	}

	// Step 4: Restore original values.
	restored := respGet.DHCPv4
	restored.General.ValidLifetime = origValidLifetime
	restored.General.FWRules = origFWRules

	respSet, err = controller.SettingsDHCPv4Update(ctx, &restored)
	if err != nil {
		t.Fatalf("SettingsDHCPv4Update (restore) failed: %v", err)
	}
	t.Logf("SettingsDHCPv4Update (restore): %+v", respSet)

	// Step 5: Verify restore.
	respGet, err = controller.SettingsDHCPv4Get(ctx)
	if err != nil {
		t.Fatalf("SettingsDHCPv4Get (after restore) failed: %v", err)
	}
	if respGet.DHCPv4.General.ValidLifetime != origValidLifetime {
		t.Fatalf("valid_lifetime not restored; got %q, want %q", respGet.DHCPv4.General.ValidLifetime, origValidLifetime)
	}
	if respGet.DHCPv4.General.FWRules != origFWRules {
		t.Fatalf("fwrules not restored; got %q, want %q", respGet.DHCPv4.General.FWRules, origFWRules)
	}
}

func TestSettingsDHCPv6(t *testing.T) {
	api_client := api.NewClient(api.Options{
		Uri:           os.Getenv("OPNSENSE_URI"),
		APIKey:        os.Getenv("OPNSENSE_API_KEY"),
		APISecret:     os.Getenv("OPNSENSE_API_SECRET"),
		AllowInsecure: true,
		MaxBackoff:    30,
		MinBackoff:    1,
		MaxRetries:    4,
	})

	controller := Controller{Api: api_client}
	ctx := context.Background()

	// Step 1: Get current settings.
	respGet, err := controller.SettingsDHCPv6Get(ctx)
	if err != nil {
		t.Fatalf("SettingsDHCPv6Get failed: %v", err)
	}
	t.Logf("SettingsDHCPv6Get (initial): %+v", respGet)

	// Verify key fields are present.
	if respGet.DHCPv6.General.ValidLifetime == "" {
		t.Fatal("Expected general.valid_lifetime to be set")
	}
	if respGet.DHCPv6.General.Enabled == "" {
		t.Fatal("Expected general.enabled to be set")
	}

	// Preserve originals.
	origValidLifetime := respGet.DHCPv6.General.ValidLifetime
	origFWRules := respGet.DHCPv6.General.FWRules

	// Step 2: Update safe fields.
	updated := respGet.DHCPv6
	updated.General.ValidLifetime = "3600"
	updated.General.FWRules = "0"

	respSet, err := controller.SettingsDHCPv6Update(ctx, &updated)
	if err != nil {
		t.Fatalf("SettingsDHCPv6Update failed: %v", err)
	}
	t.Logf("SettingsDHCPv6Update (update): %+v", respSet)
	if respSet.Result != "saved" {
		t.Fatalf("Expected result 'saved', got %q", respSet.Result)
	}

	// Step 3: Read back and verify.
	respGet, err = controller.SettingsDHCPv6Get(ctx)
	if err != nil {
		t.Fatalf("SettingsDHCPv6Get (after update) failed: %v", err)
	}
	if respGet.DHCPv6.General.ValidLifetime != "3600" {
		t.Fatalf("valid_lifetime not updated; got %q, want %q", respGet.DHCPv6.General.ValidLifetime, "3600")
	}
	if respGet.DHCPv6.General.FWRules != "0" {
		t.Fatalf("fwrules not updated; got %q, want %q", respGet.DHCPv6.General.FWRules, "0")
	}

	// Step 4: Restore.
	restored := respGet.DHCPv6
	restored.General.ValidLifetime = origValidLifetime
	restored.General.FWRules = origFWRules

	respSet, err = controller.SettingsDHCPv6Update(ctx, &restored)
	if err != nil {
		t.Fatalf("SettingsDHCPv6Update (restore) failed: %v", err)
	}
	t.Logf("SettingsDHCPv6Update (restore): %+v", respSet)

	// Step 5: Verify restore.
	respGet, err = controller.SettingsDHCPv6Get(ctx)
	if err != nil {
		t.Fatalf("SettingsDHCPv6Get (after restore) failed: %v", err)
	}
	if respGet.DHCPv6.General.ValidLifetime != origValidLifetime {
		t.Fatalf("valid_lifetime not restored; got %q, want %q", respGet.DHCPv6.General.ValidLifetime, origValidLifetime)
	}
	if respGet.DHCPv6.General.FWRules != origFWRules {
		t.Fatalf("fwrules not restored; got %q, want %q", respGet.DHCPv6.General.FWRules, origFWRules)
	}
}

func TestSettingsReconfigure(t *testing.T) {
	api_client := api.NewClient(api.Options{
		Uri:           os.Getenv("OPNSENSE_URI"),
		APIKey:        os.Getenv("OPNSENSE_API_KEY"),
		APISecret:     os.Getenv("OPNSENSE_API_SECRET"),
		AllowInsecure: true,
		MaxBackoff:    30,
		MinBackoff:    1,
		MaxRetries:    4,
	})

	controller := Controller{Api: api_client}
	ctx := context.Background()

	resp, err := controller.SettingsReconfigure(ctx)
	if err != nil {
		t.Fatalf("SettingsReconfigure failed: %v", err)
	}
	t.Logf("SettingsReconfigure: %+v", resp)
	if resp.Status != "ok" {
		t.Fatalf("Expected status 'ok', got %q", resp.Status)
	}
}
