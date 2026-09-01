package firewall

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func newController(t *testing.T) Controller {
	t.Helper()
	return Controller{
		Api: api.NewClient(api.Options{
			Uri:           os.Getenv("OPNSENSE_URI"),
			APIKey:        os.Getenv("OPNSENSE_API_KEY"),
			APISecret:     os.Getenv("OPNSENSE_API_SECRET"),
			AllowInsecure: true,
			MaxBackoff:    30,
			MinBackoff:    1,
			MaxRetries:    4,
		}),
	}
}

func TestAlias(t *testing.T) {
	controller := newController(t)
	ctx := context.Background()

	alias := &Alias{
		Enabled:     "1",
		Name:        "testalias",
		Type:        api.SelectedMap("host"),
		IPProtocol:  api.SelectedMapList([]string{"IPv4"}),
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

func TestAliasURLJSON(t *testing.T) {
	controller := newController(t)
	ctx := context.Background()

	alias := &Alias{
		Enabled:        "1",
		Name:           "testurljsonalias",
		Type:           api.SelectedMap("urljson"),
		Content:        api.SelectedMapListNL{"https://api.github.com/meta"},
		PathExpression: ".web + .api + .git | .[]",
		Description:    "Test urljson alias",
	}

	key, err := controller.AddAlias(ctx, alias)
	if err != nil {
		t.Fatalf("Failed to add urljson alias: %v", err)
	}
	t.Logf("Added urljson alias with key: %s", key)

	retrievedAlias, err := controller.GetAlias(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get urljson alias: %v", err)
	}
	t.Logf("Retrieved urljson alias: %+v", retrievedAlias)
	if retrievedAlias.Type.String() != "urljson" {
		t.Fatalf("Retrieved alias type does not match: got %s, want urljson", retrievedAlias.Type.String())
	}
	if retrievedAlias.PathExpression != alias.PathExpression {
		t.Fatalf("Retrieved alias path_expression does not match: got %s, want %s", retrievedAlias.PathExpression, alias.PathExpression)
	}

	alias.PathExpression = ".web | .[]"
	alias.Description = "Test urljson alias updated"
	err = controller.UpdateAlias(ctx, key, alias)
	if err != nil {
		t.Fatalf("Failed to update urljson alias: %v", err)
	}

	retrievedAlias, err = controller.GetAlias(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get updated urljson alias: %v", err)
	}
	if retrievedAlias.PathExpression != ".web | .[]" {
		t.Fatalf("Retrieved alias path_expression does not match updated value: got %s, want '.web | .[]'", retrievedAlias.PathExpression)
	}

	err = controller.DeleteAlias(ctx, key)
	if err != nil {
		t.Fatalf("Failed to delete urljson alias: %v", err)
	}
	t.Logf("Deleted urljson alias with key: %s", key)
}
