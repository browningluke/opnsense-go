package dnsmasq

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestSearch(t *testing.T) {
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

	respHost, err := controller.SearchHost(ctx, "-1")
	if err != nil {
		t.Fatalf("Failed to search Hosts: %v", err)
	}
	t.Logf("SearchHost: %+v", respHost)

	respDomain, err := controller.SearchDomain(ctx, "-1")
	if err != nil {
		t.Fatalf("Failed to search Domains: %v", err)
	}
	t.Logf("SearchDomain: %+v", respDomain)

	respRange, err := controller.SearchRange(ctx, "-1")
	if err != nil {
		t.Fatalf("Failed to search Ranges: %v", err)
	}
	t.Logf("respRange: %+v", respRange)

	respOption, err := controller.SearchOption(ctx, "-1")
	if err != nil {
		t.Fatalf("Failed to search Options: %v", err)
	}
	t.Logf("SearchOption: %+v", respOption)

	respBoot, err := controller.SearchBoot(ctx, "-1")
	if err != nil {
		t.Fatalf("Failed to search Boots: %v", err)
	}
	t.Logf("SearchBoot: %+v", respBoot)

	respTag, err := controller.SearchTag(ctx, "-1")
	if err != nil {
		t.Fatalf("Failed to search Tags: %v", err)
	}
	t.Logf("SearchTag: %+v", respTag)
}
