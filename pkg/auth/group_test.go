package auth

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestGroup(t *testing.T) {
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

	group := &Group{
		Name:        "test-group",
		Description: "test group description",
		Privilege:   api.SelectedMapList([]string{"page-diagnostics-authentication", "page-diagnostics-backup-restore"}),
	}

	key, err := controller.AddGroup(ctx, group)
	if err != nil {
		t.Fatalf("Failed to add group: %v", err)
	}
	t.Logf("Added group with key: %s", key)

	retrievedGroup, err := controller.GetGroup(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get group: %v", err)
	}
	t.Logf("Retrieved group: %+v", retrievedGroup)
	if retrievedGroup.Name != group.Name {
		t.Fatalf("Retrieved group name does not match: got %s, want %s", retrievedGroup.Name, group.Name)
	}
	if retrievedGroup.Description != group.Description {
		t.Fatalf("Retrieved group description does not match: got %s, want %s", retrievedGroup.Description, group.Description)
	}

	group.Name = "test-group-updated"
	group.Privilege = api.SelectedMapList([]string{"page-diagnostics-authentication"})
	group.Member = api.SelectedMap("0")
	err = controller.UpdateGroup(ctx, key, group)
	if err != nil {
		t.Fatalf("Failed to update group: %v", err)
	}
	t.Logf("Updated group with key: %s", key)

	retrievedGroup, err = controller.GetGroup(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get updated group: %v", err)
	}
	if retrievedGroup.Name != group.Name {
		t.Fatalf("Retrieved group name does not match updated name: got %s, want %s", retrievedGroup.Name, group.Name)
	}
	if retrievedGroup.Privilege.String() != group.Privilege.String() {
		t.Fatalf("Retrieved group privileges does not match updated privileges: got %s, want %s", retrievedGroup.Privilege.String(), group.Privilege.String())
	}

	err = controller.DeleteGroup(ctx, key)
	if err != nil {
		t.Fatalf("Failed to delete group: %v", err)
	}
	t.Logf("Deleted group with key: %s", key)
}
