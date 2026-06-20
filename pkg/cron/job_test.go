package cron

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestJob(t *testing.T) {
	controller := Controller{
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
	ctx := context.Background()

	job := &Job{
		Enabled:     "1",
		Minutes:     "0",
		Hours:       "4",
		Days:        "*",
		Months:      "*",
		Weekdays:    "*",
		Who:         "root",
		Command:     api.SelectedMap("firmware poll"),
		Parameters:  "",
		Description: "opnsense-go acceptance test",
	}

	id, err := controller.AddJob(ctx, job)
	if err != nil {
		t.Fatalf("AddJob failed: %v", err)
	}
	t.Logf("Added job: %s", id)

	got, err := controller.GetJob(ctx, id)
	if err != nil {
		t.Fatalf("GetJob failed: %v", err)
	}
	if got.Description != job.Description {
		t.Fatalf("Description mismatch: got %q, want %q", got.Description, job.Description)
	}
	if got.Command.String() != job.Command.String() {
		t.Fatalf("Command mismatch: got %q, want %q", got.Command.String(), job.Command.String())
	}

	job.Hours = "3"
	job.Description = "opnsense-go acceptance test (updated)"
	if err = controller.UpdateJob(ctx, id, job); err != nil {
		t.Fatalf("UpdateJob failed: %v", err)
	}

	got, err = controller.GetJob(ctx, id)
	if err != nil {
		t.Fatalf("GetJob (after update) failed: %v", err)
	}
	if got.Hours != job.Hours {
		t.Fatalf("Hours mismatch after update: got %q, want %q", got.Hours, job.Hours)
	}
	if got.Description != job.Description {
		t.Fatalf("Description mismatch after update: got %q, want %q", got.Description, job.Description)
	}

	if err = controller.DeleteJob(ctx, id); err != nil {
		t.Fatalf("DeleteJob failed: %v", err)
	}
	t.Logf("Deleted job: %s", id)
}
