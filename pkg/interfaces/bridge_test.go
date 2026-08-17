package interfaces

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestBridge(t *testing.T) {
	// OPNsense requires at least one member interface, and bridge changes are
	// applied immediately. Only run where a safe member is explicitly provided
	// (CI sets lan on its throwaway VM).
	member := os.Getenv("OPNSENSE_TEST_BRIDGE_MEMBER")
	if member == "" {
		t.Skip("OPNSENSE_TEST_BRIDGE_MEMBER not set; skipping bridge integration test")
	}

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

	bridge := &Bridge{
		Description: "Test Bridge",
		Members:     api.SelectedMapList{member},
		LinkLocal:   "0",
		EnableStp:   "0",
		StpProto:    "rstp",
	}

	key, err := controller.AddBridge(ctx, bridge)
	if err != nil {
		t.Fatalf("Failed to add bridge: %v", err)
	}
	t.Logf("Added bridge with key: %s", key)

	retrievedBridge, err := controller.GetBridge(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get bridge: %v", err)
	}
	t.Logf("Retrieved bridge: %+v", retrievedBridge)
	if retrievedBridge.Description != bridge.Description {
		t.Errorf("Retrieved bridge description does not match: got %s, want %s", retrievedBridge.Description, bridge.Description)
	}
	if retrievedBridge.EnableStp != bridge.EnableStp {
		t.Errorf("Retrieved bridge enablestp does not match: got %s, want %s", retrievedBridge.EnableStp, bridge.EnableStp)
	}
	if retrievedBridge.Device == "" {
		t.Errorf("Retrieved bridge has no device name assigned")
	}

	retrievedBridge.Description = "Test Bridge (updated)"
	err = controller.UpdateBridge(ctx, key, retrievedBridge)
	if err != nil {
		t.Fatalf("Failed to update bridge: %v", err)
	}

	updatedBridge, err := controller.GetBridge(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get updated bridge: %v", err)
	}
	if updatedBridge.Description != "Test Bridge (updated)" {
		t.Errorf("Updated bridge description does not match: got %s", updatedBridge.Description)
	}

	// Clean up the bridge after the test
	err = controller.DeleteBridge(ctx, key)
	if err != nil {
		t.Fatalf("Failed to delete bridge: %v", err)
	}
	t.Logf("Deleted bridge with key: %s", key)
}
