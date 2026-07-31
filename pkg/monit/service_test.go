package monit

import (
	"context"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestService(t *testing.T) {
	controller := newController()
	ctx := context.Background()

	service := &Service{
		Enabled:      "1",
		Name:         "test-service",
		Description:  "test-service-description",
		Type:         api.SelectedMap("process"),
		Pidfile:      "/var/run/test.pid",
		Timeout:      "300",
		Starttimeout: "30",
	}

	id, err := controller.AddService(ctx, service)
	if err != nil {
		t.Fatalf("AddService failed: %v", err)
	}
	t.Logf("AddService: uuid=%s", id)

	got, err := controller.GetService(ctx, id)
	if err != nil {
		t.Fatalf("GetService failed: %v", err)
	}
	t.Logf("GetService: %+v", got)
	if got.Name != "test-service" {
		t.Fatalf("Name mismatch: got %s", got.Name)
	}

	service.Description = "test-service-description-updated"
	err = controller.UpdateService(ctx, id, service)
	if err != nil {
		t.Fatalf("UpdateService failed: %v", err)
	}

	got, err = controller.GetService(ctx, id)
	if err != nil {
		t.Fatalf("GetService after update failed: %v", err)
	}
	if got.Description != "test-service-description-updated" {
		t.Fatalf("Description not updated: got %s", got.Description)
	}

	err = controller.DeleteService(ctx, id)
	if err != nil {
		t.Fatalf("DeleteService failed: %v", err)
	}
	t.Log("DeleteService: deleted")
}
