package firewall

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestAlias(t *testing.T) {
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

	alias := &Alias{
		Enabled:     "1",
		Name:        "testalias",
		Type:        api.SelectedMap("host"),
		IPProtocol:  api.SelectedMap("IPv4"),
		Content:     api.SelectedMapListNL{"192.168.1.1"},
		Description: "Test alias",
	}

	key, err := controller.AddAlias(ctx, alias)
	if err != nil {
		t.Fatalf("Failed to add alias: %v", err)
	}
	t.Logf("Added alias with key: %s", key)

	retrievedAlias, err := controller.GetAlias(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get alias: %v", err)
	}
	t.Logf("Retrieved alias: %+v", retrievedAlias)
	if retrievedAlias.Name != alias.Name {
		t.Fatalf("Retrieved alias name does not match: got %s, want %s", retrievedAlias.Name, alias.Name)
	}
	if retrievedAlias.Description != alias.Description {
		t.Fatalf("Retrieved alias description does not match: got %s, want %s", retrievedAlias.Description, alias.Description)
	}

	alias.Name = "testaliasupd"
	alias.Description = "Test alias updated"
	err = controller.UpdateAlias(ctx, key, alias)
	if err != nil {
		t.Fatalf("Failed to update alias: %v", err)
	}
	t.Logf("Updated alias with key: %s", key)

	retrievedAlias, err = controller.GetAlias(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get updated alias: %v", err)
	}
	if retrievedAlias.Name != "testaliasupd" {
		t.Fatalf("Retrieved alias name does not match updated name: got %s, want %s", retrievedAlias.Name, "testaliasupd")
	}
	if retrievedAlias.Description != "Test alias updated" {
		t.Fatalf("Retrieved alias description does not match updated description: got %s, want %s", retrievedAlias.Description, "Test alias updated")
	}

	err = controller.DeleteAlias(ctx, key)
	if err != nil {
		t.Fatalf("Failed to delete alias: %v", err)
	}
	t.Logf("Deleted alias with key: %s", key)
}
