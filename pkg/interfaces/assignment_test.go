package interfaces

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestAssignment(t *testing.T) {
	apiClient := api.NewClient(api.Options{
		Uri:           os.Getenv("OPNSENSE_URI"),
		APIKey:        os.Getenv("OPNSENSE_API_KEY"),
		APISecret:     os.Getenv("OPNSENSE_API_SECRET"),
		AllowInsecure: true,
		MaxBackoff:    30,
		MinBackoff:    1,
		MaxRetries:    4,
	})

	controller := Controller{Api: apiClient}
	ctx := context.Background()

	vlanID, err := controller.AddVlan(ctx, &Vlan{
		Priority:    "0",
		Description: "Assignment test VLAN",
		Parent:      "vtnet0",
		Tag:         "4093",
		Device:      "vlan04093",
	})
	if err != nil {
		t.Fatalf("Failed to add test VLAN: %v", err)
	}
	defer func() {
		if err := controller.DeleteVlan(ctx, vlanID); err != nil {
			t.Errorf("Failed to delete test VLAN: %v", err)
		}
	}()

	assignment := &Assignment{
		Device:      api.SelectedMap("vlan04093"),
		Description: "AssignmentTest",
		Lock:        "1",
	}

	assignmentID, err := controller.AddAssignment(ctx, assignment)
	if err != nil {
		// The integration workflow currently exercises OPNsense 26.1, while
		// the official assignment controller was introduced in 26.7.
		if strings.Contains(err.Error(), "status code 404") {
			t.Skip("interface assignment API requires OPNsense 26.7 or newer")
		}
		t.Fatalf("Failed to add interface assignment: %v", err)
	}
	defer func() {
		assignment.Lock = "0"
		if err := controller.UpdateAssignment(ctx, assignmentID, assignment); err != nil {
			t.Errorf("Failed to unlock interface assignment: %v", err)
			return
		}
		if err := controller.DeleteAssignment(ctx, assignmentID); err != nil {
			t.Errorf("Failed to delete interface assignment: %v", err)
		}
	}()

	retrieved, err := controller.GetAssignment(ctx, assignmentID)
	if err != nil {
		t.Fatalf("Failed to get interface assignment: %v", err)
	}
	if retrieved.Device.String() != assignment.Device.String() {
		t.Errorf("Retrieved device does not match: got %s, want %s", retrieved.Device, assignment.Device)
	}
	retrieved.Description = "AssignmentUpdated"
	if err := controller.UpdateAssignment(ctx, assignmentID, retrieved); err != nil {
		t.Fatalf("Failed to update interface assignment: %v", err)
	}

	updated, err := controller.GetAssignment(ctx, assignmentID)
	if err != nil {
		t.Fatalf("Failed to get updated interface assignment: %v", err)
	}
	if updated.Description != "AssignmentUpdated" {
		t.Errorf("Updated description does not match: got %s", updated.Description)
	}
}
