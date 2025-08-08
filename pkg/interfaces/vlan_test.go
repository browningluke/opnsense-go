package interfaces

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestVLAN(t *testing.T) {
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

	vlan := &Vlan{
		Priority:    "1",
		Description: "Test VLAN",
		Parent:      "vtnet0",
		Tag:         "100",
		Device:      "", // Leave empty for auto-assignment
	}

	key, err := controller.AddVlan(ctx, vlan)
	if err != nil {
		t.Fatalf("Failed to add VLAN: %v", err)
	}
	t.Logf("Added VLAN with key: %s", key)

	retrievedVLAN, err := controller.GetVlan(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get VLAN: %v", err)
	}
	t.Logf("Retrieved VLAN: %+v", retrievedVLAN)
	if retrievedVLAN.Priority != vlan.Priority {
		t.Errorf("Retrieved VLAN priority does not match original VLAN")
	}
	if retrievedVLAN.Parent != vlan.Parent {
		t.Errorf("Retrieved VLAN parent does not match original VLAN parent: got %s, want %s", retrievedVLAN.Parent, vlan.Parent)
	}
	if retrievedVLAN.Tag != vlan.Tag {
		t.Errorf("Retrieved VLAN tag does not match original VLAN tag: got %s, want %s", retrievedVLAN.Tag, vlan.Tag)
	}
	if retrievedVLAN.Description != vlan.Description {
		t.Errorf("Retrieved VLAN description does not match original VLAN description: got %s, want %s", retrievedVLAN.Description, vlan.Description)
	}

	// Clean up the VLAN after the test
	err = controller.DeleteVlan(ctx, key)
	if err != nil {
		t.Fatalf("Failed to delete VLAN: %v", err)
	}
	t.Logf("Deleted VLAN with key: %s", key)
}
