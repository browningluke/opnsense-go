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

	rng := &Range{
		StartAddress: "192.168.100.100",
		EndAddress:   "192.168.100.199",
		DomainType:   api.SelectedMap("range"),
	}

	respAdd, err := controller.AddRange(ctx, rng)
	if err != nil {
		t.Fatalf("Failed to add range: %v", err)
	}
	t.Logf("AddRange: %+v", respAdd)

	respGet, err := controller.GetRange(ctx, respAdd)
	if err != nil {
		t.Fatalf("Failed to get range: %v", err)
	}
	t.Logf("GetRange: %+v", respGet)

	rng.Interface = api.SelectedMap("lan")
	rng.Tags = api.SelectedMap("d594fa8a-1a76-44e7-afab-adb6c5bdb69e")
	rng.StartAddress = "192.168.100.200"
	rng.EndAddress = "192.168.100.250"
	rng.SubnetMask = "255.255.255.224"
	// rng.Constructor = api.SelectedMap("lan")
	// rng.Mode = api.SelectedMap("static")
	rng.LeaseTime = "3600"
	rng.Domain = "test-domain-updated"
	rng.NoSync = "1"
	rng.Description = "test-description-updated"
	err = controller.UpdateRange(ctx, respAdd, rng)
	if err != nil {
		t.Fatalf("Failed to update range: %v", err)
	}
	t.Logf("UpdateRange: %+v", rng)

	respGet, err = controller.GetRange(ctx, respAdd)
	if err != nil {
		t.Fatalf("Failed to get range: %v", err)
	}
	t.Logf("GetRange: %+v", respGet)

	err = controller.DeleteRange(ctx, respAdd)
	if err != nil {
		t.Fatalf("Failed to delete range: %v", err)
	}
	t.Log("DeleteRange: Deleted!")
}
