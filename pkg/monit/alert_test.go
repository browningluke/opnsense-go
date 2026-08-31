package monit

import (
	"context"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestAlert(t *testing.T) {
	controller := newController()
	ctx := context.Background()

	alert := &Alert{
		Enabled:     "1",
		Recipient:   "test@example.com",
		Noton:       "0",
		Events:      api.SelectedMapList{},
		Description: "test-alert",
	}

	id, err := controller.AddAlert(ctx, alert)
	if err != nil {
		t.Fatalf("AddAlert failed: %v", err)
	}
	t.Logf("AddAlert: uuid=%s", id)

	got, err := controller.GetAlert(ctx, id)
	if err != nil {
		t.Fatalf("GetAlert failed: %v", err)
	}
	t.Logf("GetAlert: %+v", got)
	if got.Recipient != "test@example.com" {
		t.Fatalf("Recipient mismatch: got %s", got.Recipient)
	}

	alert.Description = "test-alert-updated"
	err = controller.UpdateAlert(ctx, id, alert)
	if err != nil {
		t.Fatalf("UpdateAlert failed: %v", err)
	}

	got, err = controller.GetAlert(ctx, id)
	if err != nil {
		t.Fatalf("GetAlert after update failed: %v", err)
	}
	if got.Description != "test-alert-updated" {
		t.Fatalf("Description not updated: got %s", got.Description)
	}

	err = controller.DeleteAlert(ctx, id)
	if err != nil {
		t.Fatalf("DeleteAlert failed: %v", err)
	}
	t.Log("DeleteAlert: deleted")
}
