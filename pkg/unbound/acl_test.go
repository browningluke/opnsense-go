package unbound

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestAcl(t *testing.T) {
	opnsenseURL := os.Getenv("OPNSENSE_URI")
	opnsenseKey := os.Getenv("OPNSENSE_API_KEY")
	opnsenseSecret := os.Getenv("OPNSENSE_API_SECRET")

	apiClient := api.NewClient(api.Options{
		Uri:           opnsenseURL,
		APIKey:        opnsenseKey,
		APISecret:     opnsenseSecret,
		AllowInsecure: true,
		MaxBackoff:    30,
		MinBackoff:    1,
		MaxRetries:    4,
	})

	controller := Controller{Api: apiClient}
	ctx := context.Background()

	// Step 1: Create ACL.
	id, err := controller.AddAcl(ctx, &Acl{
		Enabled:     "1",
		Name:        "tf-test-acl",
		Action:      api.SelectedMap("allow"),
		Networks:    api.SelectedMapList([]string{"10.0.0.0/24", "10.0.1.0/24"}),
		Description: "integration test",
	})
	if err != nil {
		t.Fatalf("AddAcl: %v", err)
	}
	if id == "" {
		t.Fatal("AddAcl: empty id")
	}
	t.Logf("AddAcl: id=%s", id)

	// Step 2: Read back and verify.
	acl, err := controller.GetAcl(ctx, id)
	if err != nil {
		t.Fatalf("GetAcl: %v", err)
	}
	if acl.Name != "tf-test-acl" {
		t.Errorf("Name: got %q, want %q", acl.Name, "tf-test-acl")
	}
	if acl.Action.String() != "allow" {
		t.Errorf("Action: got %q, want %q", acl.Action.String(), "allow")
	}
	if acl.Enabled != "1" {
		t.Errorf("Enabled: got %q, want %q", acl.Enabled, "1")
	}
	networks := acl.Networks
	if len(networks) != 2 {
		t.Errorf("Networks: got %d entries, want 2", len(networks))
	}
	t.Logf("GetAcl: %+v", acl)

	// Step 3: Update.
	err = controller.UpdateAcl(ctx, id, &Acl{
		Enabled:     "0",
		Name:        "tf-test-acl-updated",
		Action:      api.SelectedMap("refuse"),
		Networks:    api.SelectedMapList([]string{"192.168.1.0/24"}),
		Description: "updated",
	})
	if err != nil {
		t.Fatalf("UpdateAcl: %v", err)
	}

	// Step 4: Read back and verify update.
	acl, err = controller.GetAcl(ctx, id)
	if err != nil {
		t.Fatalf("GetAcl after update: %v", err)
	}
	if acl.Name != "tf-test-acl-updated" {
		t.Errorf("Name after update: got %q, want %q", acl.Name, "tf-test-acl-updated")
	}
	if acl.Action.String() != "refuse" {
		t.Errorf("Action after update: got %q, want %q", acl.Action.String(), "refuse")
	}
	if acl.Enabled != "0" {
		t.Errorf("Enabled after update: got %q, want %q", acl.Enabled, "0")
	}
	t.Logf("GetAcl after update: %+v", acl)

	// Step 5: Delete.
	err = controller.DeleteAcl(ctx, id)
	if err != nil {
		t.Fatalf("DeleteAcl: %v", err)
	}

	// Step 6: Verify deletion.
	_, err = controller.GetAcl(ctx, id)
	if err == nil {
		t.Fatal("GetAcl after delete: expected error, got nil")
	}
	t.Logf("GetAcl after delete: %v (expected error)", err)
}
