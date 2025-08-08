package firewall

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestFilter(t *testing.T) {
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

	ctx := context.Background()

	controller := Controller{
		Api: api_client,
	}

	filter := &Filter{
		Enabled:           "1",
		Sequence:          "1",
		Action:            api.SelectedMap("pass"),
		Quick:             "1",
		Interface:         api.SelectedMapList{"wan"},
		Direction:         api.SelectedMap("in"),
		IPProtocol:        api.SelectedMap("inet"),
		Protocol:          api.SelectedMap("TCP"),
		SourceNet:         "any",
		SourcePort:        "",
		SourceInvert:      "0",
		DestinationNet:    "192.168.1.0/24",
		DestinationPort:   "80",
		DestinationInvert: "0",
		Log:               "0",
		Description:       "Test filter rule",
	}

	key, err := controller.AddFilter(ctx, filter)
	if err != nil {
		t.Fatalf("Failed to add filter rule: %v", err)
	}
	t.Logf("Added filter rule with key: %s", key)

	retrievedFilter, err := controller.GetFilter(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get filter rule: %v", err)
	}
	t.Logf("Retrieved filter rule: %+v", retrievedFilter)
	if retrievedFilter.SourceNet != filter.SourceNet {
		t.Fatalf("Retrieved filter source net does not match: got %s, want %s", retrievedFilter.SourceNet, filter.SourceNet)
	}
	if retrievedFilter.DestinationNet != filter.DestinationNet {
		t.Fatalf("Retrieved filter destination net does not match: got %s, want %s", retrievedFilter.DestinationNet, filter.DestinationNet)
	}
	if retrievedFilter.Description != filter.Description {
		t.Fatalf("Retrieved filter description does not match: got %s, want %s", retrievedFilter.Description, filter.Description)
	}

	filter.DestinationNet = "192.168.2.0/24"
	filter.DestinationPort = "443"
	filter.Description = "Test filter rule updated"
	err = controller.UpdateFilter(ctx, key, filter)
	if err != nil {
		t.Fatalf("Failed to update filter rule: %v", err)
	}
	t.Logf("Updated filter rule with key: %s", key)

	retrievedFilter, err = controller.GetFilter(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get updated filter rule: %v", err)
	}
	if retrievedFilter.DestinationNet != "192.168.2.0/24" {
		t.Fatalf("Retrieved filter destination net does not match updated net: got %s, want %s", retrievedFilter.DestinationNet, "192.168.2.0/24")
	}
	if retrievedFilter.DestinationPort != "443" {
		t.Fatalf("Retrieved filter destination port does not match updated port: got %s, want %s", retrievedFilter.DestinationPort, "443")
	}
	if retrievedFilter.Description != "Test filter rule updated" {
		t.Fatalf("Retrieved filter description does not match updated description: got %s, want %s", retrievedFilter.Description, "Test filter rule updated")
	}

	err = controller.DeleteFilter(ctx, key)
	if err != nil {
		t.Fatalf("Failed to delete filter rule: %v", err)
	}
	t.Logf("Deleted filter rule with key: %s", key)
}
