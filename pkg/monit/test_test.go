package monit

import (
	"context"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestTest(t *testing.T) {
	controller := newController()
	ctx := context.Background()

	test := &Test{
		Name:      "test-test",
		Type:      api.SelectedMap("ProgramStatus"),
		Condition: "status != 0",
		Action:    api.SelectedMap("alert"),
	}

	id, err := controller.AddTest(ctx, test)
	if err != nil {
		t.Fatalf("AddTest failed: %v", err)
	}
	t.Logf("AddTest: uuid=%s", id)

	got, err := controller.GetTest(ctx, id)
	if err != nil {
		t.Fatalf("GetTest failed: %v", err)
	}
	t.Logf("GetTest: %+v", got)
	if got.Name != "test-test" {
		t.Fatalf("Name mismatch: got %s", got.Name)
	}

	test.Condition = "status != 1"
	err = controller.UpdateTest(ctx, id, test)
	if err != nil {
		t.Fatalf("UpdateTest failed: %v", err)
	}

	got, err = controller.GetTest(ctx, id)
	if err != nil {
		t.Fatalf("GetTest after update failed: %v", err)
	}
	if got.Condition != "status != 1" {
		t.Fatalf("Condition not updated: got %s", got.Condition)
	}

	err = controller.DeleteTest(ctx, id)
	if err != nil {
		t.Fatalf("DeleteTest failed: %v", err)
	}
	t.Log("DeleteTest: deleted")
}
