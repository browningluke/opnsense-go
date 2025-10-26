package dnsmasq

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestDomain(t *testing.T) {
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

	domain := &Domain{
		Sequence: "1",
		Domain:   "test-domain",
	}

	respAdd, err := controller.AddDomain(ctx, domain)
	if err != nil {
		t.Fatalf("Failed to add Domain: %v", err)
	}
	t.Logf("AddDomain: %+v", respAdd)

	respGet, err := controller.GetDomain(ctx, respAdd)
	if err != nil {
		t.Fatalf("Failed to get Domain: %v", err)
	}
	t.Logf("GetDomain: %+v", respGet)

	domain.Ip = "192.168.100.100"
	domain.Port = "3000"
	domain.SourceIp = "0.0.0.0"
	// domain.FirewallAlias = api.SelectedMap("6a852649-aa0b-48b6-9079-3ab386b844ac")
	domain.Description = "test-description-updated"
	err = controller.UpdateDomain(ctx, respAdd, domain)
	if err != nil {
		t.Fatalf("Failed to update Domain: %v", err)
	}
	t.Logf("UpdateDomain: %+v", domain)

	respGet, err = controller.GetDomain(ctx, respAdd)
	if err != nil {
		t.Fatalf("Failed to get Domain: %v", err)
	}
	t.Logf("GetDomain: %+v", respGet)

	err = controller.DeleteDomain(ctx, respAdd)
	if err != nil {
		t.Fatalf("Failed to delete Domain: %v", err)
	}
	t.Log("DeleteDomain: Deleted!")
}
