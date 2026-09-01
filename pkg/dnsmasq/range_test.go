package dnsmasq

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestRange(t *testing.T) {
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

	expectedMode := api.SelectedMapList([]string{
		"static",
	})
	rng := &Range{
		StartAddress: "192.168.100.100",
		DomainType:   api.SelectedMap("range"),
		Mode:         expectedMode,
		Description:  "test-description",
	}

	respId, err := controller.AddRange(ctx, rng)
	if err != nil {
		t.Fatalf("Failed to add range: %v", err)
	}
	t.Logf("AddRange: %+v", respId)

	// Cleanup is only used if the test fails before the explicit delete.
	deleted := false

	t.Cleanup(func() {
		if deleted {
			return
		}

		if err := controller.DeleteRange(ctx, respId); err != nil {
			t.Errorf("Failed to cleanup range: %v", err)
		}
	})

	v, err := controller.GetRange(ctx, respId)
	if err != nil {
		t.Fatalf("Failed to get range: %v", err)
	}
	t.Logf("GetRange: %+v", v)

	if v.StartAddress != rng.StartAddress {
		t.Fatalf("StartAddress not equal; Got: %q Expected: %q", v.StartAddress, rng.StartAddress)
	}

	if v.EndAddress != "" {
		t.Fatalf("EndAddress must be empty for static range; Got: %q", v.EndAddress)
	}

	if v.SubnetMask != "" {
		t.Fatalf("SubnetMask must be empty for static range; Got: %q", v.SubnetMask)
	}

	if v.DomainType.String() != rng.DomainType.String() {
		t.Fatalf("DomainType not equal; Got: %q Expected: %q", v.DomainType.String(), rng.DomainType.String())
	}

	if v.Description != rng.Description {
		t.Fatalf("Description not equal; Got: %q Expected: %q", v.Description, rng.Description)
	}

	// ---------------------------------------------------------------------
	// Verify SelectedMapList round-trip.
	// ---------------------------------------------------------------------

	if len(v.Mode) != 1 {
		t.Fatalf("Mode contains unexpected number of values; Got: %+v Expected: %+v", v.Mode, expectedMode)
	}

	if v.Mode[0] != "static" {
		t.Fatalf("Mode not equal; Got: %+v Expected: %+v", v.Mode, expectedMode)
	}

	// ---------------------------------------------------------------------
	// Update static range.
	//
	// Do NOT specify EndAddress or SubnetMask.
	// ---------------------------------------------------------------------
	rng.StartAddress = "192.168.100.101"
	rng.Description = "test-description-updated"
	rng.Mode = api.SelectedMapList([]string{"static"})
	err = controller.UpdateRange(ctx, respId, rng)
	if err != nil {
		t.Fatalf("Failed to update range: %v", err)
	}
	t.Logf("UpdateRange: %+v", rng)

	// ---------------------------------------------------------------------
	// GET after update
	// ---------------------------------------------------------------------
	v, err = controller.GetRange(ctx, respId)
	if err != nil {
		t.Fatalf("Failed to get range after update: %v", err)
	}
	t.Logf("GetRange after update: %+v", v)

	if v.StartAddress != rng.StartAddress {
		t.Fatalf("StartAddress not updated; Got: %q Expected: %q", v.StartAddress, rng.StartAddress)
	}
	if v.EndAddress != "" {
		t.Fatalf("EndAddress must remain empty for static range; Got: %q", v.EndAddress)
	}
	if v.SubnetMask != "" {
		t.Fatalf("SubnetMask must remain empty for static range; Got: %q", v.SubnetMask)
	}
	if v.Description != rng.Description {
		t.Fatalf("Description not updated; Got: %q Expected: %q", v.Description, rng.Description)
	}
	if len(v.Mode) != 1 {
		t.Fatalf("Mode contains unexpected number of values after update; Got: %+v Expected: %+v", v.Mode, rng.Mode)
	}
	if v.Mode[0] != "static" {
		t.Fatalf("Mode not updated; Got: %+v Expected: %+v", v.Mode, rng.Mode)
	}

	// ---------------------------------------------------------------------
	// Delete
	// ---------------------------------------------------------------------

	err = controller.DeleteRange(ctx, respId)
	if err != nil {
		t.Fatalf("Failed to delete range: %v", err)
	}
	deleted = true
	t.Log("DeleteRange: Deleted!")
}
