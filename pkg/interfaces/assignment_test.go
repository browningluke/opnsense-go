package interfaces

import (
	"context"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

// assignmentAPIMinVersion is the first OPNsense release with AssignmentController
// (/interfaces/assignment/*); it doesn't exist on 26.1, which the CI VM image
// is pinned to.
var assignmentAPIMinVersion = [2]int{26, 7}

// skipUnlessAssignmentAPISupported skips the calling test unless the target
// OPNsense version (read from OPNSENSE_VERSION, set by the CI workflow) is
// known to be at or above assignmentAPIMinVersion. Falls back to requiring
// OPNSENSE_TEST_ASSIGNMENT_API=1 when OPNSENSE_VERSION isn't set, e.g. for
// manual runs against real hardware.
func skipUnlessAssignmentAPISupported(t *testing.T) {
	t.Helper()

	version := os.Getenv("OPNSENSE_VERSION")
	if version == "" {
		if os.Getenv("OPNSENSE_TEST_ASSIGNMENT_API") != "1" {
			t.Skip("OPNSENSE_TEST_ASSIGNMENT_API=1 required: AssignmentController endpoint only exists on OPNsense 26.7+")
		}
		return
	}

	major, minor, ok := parseMajorMinor(version)
	if !ok || major < assignmentAPIMinVersion[0] ||
		(major == assignmentAPIMinVersion[0] && minor < assignmentAPIMinVersion[1]) {
		t.Skipf("AssignmentController endpoint only exists on OPNsense 26.7+, target is %s", version)
	}
}

func parseMajorMinor(version string) (major, minor int, ok bool) {
	parts := strings.SplitN(version, ".", 3)
	if len(parts) < 2 {
		return 0, 0, false
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, false
	}
	minor, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, false
	}
	return major, minor, true
}

// TestAssignmentUpdate exercises Get/Update against the pre-existing "wan"
// assignment, since every OPNsense install already has one. It only ever
// touches Description/Lock, never Device: the CI VM this runs against has a
// single NIC assigned to wan, and repointing it would cut off the VM's only
// network path.
func TestAssignmentUpdate(t *testing.T) {
	skipUnlessAssignmentAPISupported(t)

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

	original, err := controller.GetAssignment(ctx, "wan")
	if err != nil {
		t.Fatalf("Failed to get wan assignment: %v", err)
	}
	t.Logf("Original wan assignment: %+v", original)

	updated := &Assignment{
		Device:      original.Device,
		Description: "Test WAN description",
		Lock:        original.Lock,
	}

	if err := controller.UpdateAssignment(ctx, "wan", updated); err != nil {
		t.Fatalf("Failed to update wan assignment: %v", err)
	}

	// Restore the original description so the test doesn't leave the VM
	// (or a real box, when run manually) in a modified state.
	defer func() {
		if err := controller.UpdateAssignment(ctx, "wan", original); err != nil {
			t.Errorf("Failed to restore original wan assignment: %v", err)
		}
	}()

	retrieved, err := controller.GetAssignment(ctx, "wan")
	if err != nil {
		t.Fatalf("Failed to get updated wan assignment: %v", err)
	}
	if retrieved.Description != updated.Description {
		t.Errorf("Retrieved wan description does not match: got %s, want %s", retrieved.Description, updated.Description)
	}
	if retrieved.Device != original.Device {
		t.Errorf("wan device changed unexpectedly: got %s, want %s", retrieved.Device, original.Device)
	}
}

// TestAssignmentCreateDelete exercises the full Add/Delete lifecycle against
// a spare physical device. Skipped unless OPNSENSE_TEST_ASSIGNMENT_SPARE_DEVICE
// is set to a device name that is safe to assign and unassign (i.e. not
// carrying live traffic) — the CI VM is single-NIC and has no such device.
// Also requires OPNsense 26.7+ (see skipUnlessAssignmentAPISupported).
func TestAssignmentCreateDelete(t *testing.T) {
	skipUnlessAssignmentAPISupported(t)

	device := os.Getenv("OPNSENSE_TEST_ASSIGNMENT_SPARE_DEVICE")
	if device == "" {
		t.Skip("OPNSENSE_TEST_ASSIGNMENT_SPARE_DEVICE must be set to a spare, non-live device to run this test")
	}

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

	assignment := &Assignment{
		Device:      device,
		Description: "Test assignment",
		Lock:        "0",
	}

	key, err := controller.AddAssignment(ctx, assignment)
	if err != nil {
		t.Fatalf("Failed to add assignment: %v", err)
	}
	t.Logf("Added assignment with key: %s", key)

	retrieved, err := controller.GetAssignment(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get assignment: %v", err)
	}
	if retrieved.Device != assignment.Device {
		t.Errorf("Retrieved assignment device does not match: got %s, want %s", retrieved.Device, assignment.Device)
	}
	if retrieved.Description != assignment.Description {
		t.Errorf("Retrieved assignment description does not match: got %s, want %s", retrieved.Description, assignment.Description)
	}

	if err := controller.DeleteAssignment(ctx, key); err != nil {
		t.Fatalf("Failed to delete assignment: %v", err)
	}
	t.Logf("Deleted assignment with key: %s", key)
}
