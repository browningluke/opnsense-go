package trust

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func newController() *Controller {
	return &Controller{
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

func TestCa(t *testing.T) {
	controller := newController()
	ctx := context.Background()

	ca := &Ca{
		Description: "test-ca",
		Action:      api.SelectedMap("internal"),
		KeyType:     api.SelectedMap("2048"),
		Digest:      api.SelectedMap("sha256"),
		Lifetime:    "3650",
		Country:     api.SelectedMap("US"),
		CommonName:  "Test Internal CA",
	}

	id, err := controller.AddCa(ctx, ca)
	if err != nil {
		t.Fatalf("AddCa failed: %v", err)
	}
	t.Logf("AddCa: uuid=%s", id)

	got, err := controller.GetCa(ctx, id)
	if err != nil {
		t.Fatalf("GetCa failed: %v", err)
	}
	t.Logf("GetCa: %+v", got)
	if got.CommonName != "Test Internal CA" {
		t.Fatalf("CommonName mismatch: got %s", got.CommonName)
	}
	if got.RefId == "" {
		t.Fatal("RefId should be populated after creation")
	}
	if got.CrtPayload == "" {
		t.Fatal("CrtPayload should be populated after creation")
	}

	ca.Description = "test-ca-updated"
	err = controller.UpdateCa(ctx, id, ca)
	if err != nil {
		t.Fatalf("UpdateCa failed: %v", err)
	}

	got, err = controller.GetCa(ctx, id)
	if err != nil {
		t.Fatalf("GetCa after update failed: %v", err)
	}
	if got.Description != "test-ca-updated" {
		t.Fatalf("Description not updated: got %s", got.Description)
	}

	err = controller.DeleteCa(ctx, id)
	if err != nil {
		t.Fatalf("DeleteCa failed: %v", err)
	}
	t.Log("DeleteCa: deleted")
}

func TestSettings(t *testing.T) {
	controller := newController()
	ctx := context.Background()

	result, err := controller.SettingsGet(ctx)
	if err != nil {
		t.Fatalf("SettingsGet failed: %v", err)
	}
	t.Logf("SettingsGet: %+v", result.Trust)

	// Update settings (no-op, set same values back)
	_, err = controller.SettingsSet(ctx, &result.Trust)
	if err != nil {
		t.Fatalf("SettingsSet failed: %v", err)
	}

	_, err = controller.SettingsReconfigure(ctx)
	if err != nil {
		t.Fatalf("SettingsReconfigure failed: %v", err)
	}
	t.Log("SettingsReconfigure: ok")
}
