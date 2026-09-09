package firewall

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestGroup(t *testing.T) {
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

	controller := Controller{
		Api: apiClient,
	}

	ctx := context.Background()

	group := &Group{
		Name: "test_group",
		Members: api.SelectedMapList{
			"wan",
		},
		Sequence:    "100",
		NoGroup:     "0",
		Description: "test-description",
	}

	// Add
	respAdd, err := controller.AddGroup(ctx, group)
	if err != nil {
		t.Fatalf("Failed to add group: %v", err)
	}
	t.Logf("AddGroup: %s", respAdd)

	// Auto cleanup
	t.Cleanup(func() {
		if err := controller.DeleteGroup(ctx, respAdd); err != nil {
			t.Errorf("Failed to cleanup group %s: %v", respAdd, err)
		} else {
			t.Logf("Cleanup: DeleteGroup %s", respAdd)
		}
	})

	// Get
	respGet, err := controller.GetGroup(ctx, respAdd)
	if err != nil {
		t.Fatalf("Failed to get group: %v", err)
	}
	t.Logf("GetGroup: %+v", respGet)

	if respGet.Name != group.Name {
		t.Fatalf("Name mismatch; Got: %s Expected: %s", respGet.Name, group.Name)
	}

	if respGet.Sequence != group.Sequence {
		t.Fatalf("Sequence mismatch; Got: %s Expected: %s", respGet.Sequence, group.Sequence)
	}

	if respGet.NoGroup != group.NoGroup {
		t.Fatalf("NoGroup mismatch; Got: %s Expected: %s", respGet.NoGroup, group.NoGroup)
	}

	if respGet.Description != group.Description {
		t.Fatalf("Description mismatch; Got: %s Expected: %s", respGet.Description, group.Description)
	}

	if len(respGet.Members) != len(group.Members) {
		t.Fatalf("Members count mismatch; Got: %d Expected: %d", len(respGet.Members), len(group.Members))
	}

	// Update
	group.Name = "test_group_upda" // max. 15 chars
	group.Sequence = "200"
	group.NoGroup = "1"
	group.Description = "test-description-updated"

	err = controller.UpdateGroup(ctx, respAdd, group)
	if err != nil {
		t.Fatalf("Failed to update group: %v", err)
	}
	t.Logf("UpdateGroup: %+v", group)

	// Get after update
	respGet, err = controller.GetGroup(ctx, respAdd)
	if err != nil {
		t.Fatalf("Failed to get updated group: %v", err)
	}
	t.Logf("GetGroup after update: %+v", respGet)

	if respGet.Name != group.Name {
		t.Fatalf("Name not updated; Got: %s Expected: %s", respGet.Name, group.Name)
	}

	if respGet.Sequence != group.Sequence {
		t.Fatalf("Sequence not updated; Got: %s Expected: %s", respGet.Sequence, group.Sequence)
	}

	if respGet.NoGroup != group.NoGroup {
		t.Fatalf("NoGroup not updated; Got: %s Expected: %s", respGet.NoGroup, group.NoGroup)
	}

	if respGet.Description != group.Description {
		t.Fatalf("Description not updated; Got: %s Expected: %s", respGet.Description, group.Description)
	}
}
