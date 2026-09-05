package ids

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func testController() *Controller {
	return &Controller{Api: api.NewClient(api.Options{
		Uri:           os.Getenv("OPNSENSE_URI"),
		APIKey:        os.Getenv("OPNSENSE_API_KEY"),
		APISecret:     os.Getenv("OPNSENSE_API_SECRET"),
		AllowInsecure: os.Getenv("OPNSENSE_ALLOW_INSECURE") == "true",
	})}
}

func TestSettingsGet(t *testing.T) {
	settings, err := testController().SettingsGet(context.Background())
	if err != nil {
		t.Fatalf("SettingsGet failed: %v", err)
	}
	if settings.Ids.General.Enabled == "" {
		t.Fatal("SettingsGet returned an empty general.enabled value")
	}
}

func TestPolicyLifecycle(t *testing.T) {
	ctx := context.Background()
	controller := testController()
	policy := &Policy{
		Enabled:     "0",
		Priority:    "1000",
		NewAction:   api.SelectedMap("alert"),
		Description: "opnsense-go acceptance policy",
	}

	id, err := controller.AddPolicy(ctx, policy)
	if err != nil {
		t.Fatalf("AddPolicy failed: %v", err)
	}
	t.Cleanup(func() {
		if err := controller.DeletePolicy(context.Background(), id); err != nil {
			t.Errorf("DeletePolicy cleanup failed: %v", err)
		}
	})

	created, err := controller.GetPolicy(ctx, id)
	if err != nil {
		t.Fatalf("GetPolicy failed: %v", err)
	}
	if created.Description != policy.Description {
		t.Fatalf("GetPolicy description = %q, want %q", created.Description, policy.Description)
	}

	created.Description = "opnsense-go acceptance policy updated"
	if err := controller.UpdatePolicy(ctx, id, created); err != nil {
		t.Fatalf("UpdatePolicy failed: %v", err)
	}
	updated, err := controller.GetPolicy(ctx, id)
	if err != nil {
		t.Fatalf("GetPolicy after update failed: %v", err)
	}
	if updated.Description != created.Description {
		t.Fatalf("updated description = %q, want %q", updated.Description, created.Description)
	}
}

func TestPolicyRuleLifecycle(t *testing.T) {
	ctx := context.Background()
	controller := testController()
	rule := &PolicyRule{SID: "9001001", Enabled: "0", Action: api.SelectedMap("alert"), Message: "opnsense-go acceptance override", Source: "acceptance"}

	id, err := controller.AddPolicyRule(ctx, rule)
	if err != nil {
		t.Fatalf("AddPolicyRule failed: %v", err)
	}
	t.Cleanup(func() {
		if err := controller.DeletePolicyRule(context.Background(), id); err != nil {
			t.Errorf("DeletePolicyRule cleanup failed: %v", err)
		}
	})

	created, err := controller.GetPolicyRule(ctx, id)
	if err != nil {
		t.Fatalf("GetPolicyRule failed: %v", err)
	}
	created.Message = "opnsense-go acceptance override updated"
	if err := controller.UpdatePolicyRule(ctx, id, created); err != nil {
		t.Fatalf("UpdatePolicyRule failed: %v", err)
	}
}

func TestUserRuleLifecycle(t *testing.T) {
	ctx := context.Background()
	controller := testController()
	rule := &UserRule{
		Enabled:     "0",
		Source:      "192.0.2.1",
		Destination: "any",
		Fingerprint: "CC:CC:CC:CC:CC:CC:CC:CC:CC:CC:CC:CC:CC:CC:CC:CC:CC:CC:CC:CC",
		Description: "opnsense-go acceptance user rule",
		Action:      api.SelectedMap("alert"),
		Bypass:      "0",
	}

	id, err := controller.AddUserRule(ctx, rule)
	if err != nil {
		t.Fatalf("AddUserRule failed: %v", err)
	}
	t.Cleanup(func() {
		if err := controller.DeleteUserRule(context.Background(), id); err != nil {
			t.Errorf("DeleteUserRule cleanup failed: %v", err)
		}
	})

	created, err := controller.GetUserRule(ctx, id)
	if err != nil {
		t.Fatalf("GetUserRule failed: %v", err)
	}
	created.Description = "opnsense-go acceptance user rule updated"
	if err := controller.UpdateUserRule(ctx, id, created); err != nil {
		t.Fatalf("UpdateUserRule failed: %v", err)
	}
}
