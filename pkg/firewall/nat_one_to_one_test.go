package firewall

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestNatOneToOne(t *testing.T) {
	opnsenseURL := os.Getenv("OPNSENSE_URI")
	opnsenseKey := os.Getenv("OPNSENSE_API_KEY")
	opnsenseSecret := os.Getenv("OPNSENSE_API_SECRET")

	apiClient := api.NewClient(api.Options{
		Uri:           opnsenseURL,
		APIKey:        opnsenseKey,
		APISecret:     opnsenseSecret,
		AllowInsecure: true,
		MaxBackoff:    30,
		MinBackoff:    1,
		MaxRetries:    4,
	})

	controller := Controller{Api: apiClient}
	ctx := context.Background()

	// Neues 1:1 NAT-Objekt
	rule := &NatOneToOne{
		Enabled:           "1",
		Log:               "0",
		Sequence:          "1",
		Interface:         api.SelectedMap("wan"),
		Type:              api.SelectedMap("binat"), // example "binat" or "nat" depends on your OPNsense
		SourceNet:         "192.168.2.100",          // example internal IP
		SourceInvert:      "0",
		DestinationNet:    "any",
		DestinationInvert: "0",
		ExternalNet:       "192.168.1.100", // example external IP
		NatReflection:     api.SelectedMap("enable"),
		Categories:        api.SelectedMapList{}, // no categories for now
		Description:       "Test 1:1 NAT rule",
	}

	// CREATE
	key, err := controller.AddNatOneToOne(ctx, rule)
	if err != nil {
		t.Fatalf("Failed to add NatOneToOne rule: %v", err)
	}
	t.Logf("Added NatOneToOne rule with key: %s", key)

	// READ
	got, err := controller.GetNatOneToOne(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get NatOneToOne rule: %v", err)
	}
	if got.ExternalNet != rule.ExternalNet {
		t.Errorf("ExternalNet mismatch: got %s, want %s", got.ExternalNet, rule.ExternalNet)
	}
	if got.Description != rule.Description {
		t.Errorf("Description mismatch: got %s, want %s", got.Description, rule.Description)
	}

	// UPDATE
	rule.ExternalNet = "192.168.1.101" // change external IP
	rule.SourceNet = "192.168.2.101"
	rule.Description = "Updated 1:1 NAT rule"
	err = controller.UpdateNatOneToOne(ctx, key, rule)
	if err != nil {
		t.Fatalf("Failed to update NatOneToOne rule: %v", err)
	}

	// READ after update
	updated, err := controller.GetNatOneToOne(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get updated NatOneToOne rule: %v", err)
	}
	if updated.ExternalNet != "192.168.1.101" {
		t.Errorf("ExternalNet mismatch after update: got %s, want %s", updated.ExternalNet, "192.168.1.101")
	}
	if updated.SourceNet != "192.168.2.101" {
		t.Errorf("ExternalNet mismatch after update: got %s, want %s", updated.ExternalNet, "192.168.2.101")
	}
	if updated.Description != "Updated 1:1 NAT rule" {
		t.Errorf("Description mismatch after update: got %s, want %s", updated.Description, "Updated 1:1 NAT rule")
	}

	// DELETE
	err = controller.DeleteNatOneToOne(ctx, key)
	if err != nil {
		t.Fatalf("Failed to delete NatOneToOne rule: %v", err)
	}
	t.Logf("Deleted NatOneToOne rule with key: %s", key)
}
