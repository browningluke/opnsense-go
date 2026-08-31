package syslog

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestDestination(t *testing.T) {
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

	dest := &Destination{
		Enabled:     "1",
		Transport:   api.SelectedMap("udp4"),
		Level:       api.SelectedMap("info"),
		Facility:    api.SelectedMapList{},
		Program:     api.SelectedMapList{},
		Hostname:    "192.0.2.1",
		Certificate: api.SelectedMap(""),
		Port:        "514",
		Rfc5424:     "0",
		Description: "opnsense-go acceptance test",
	}

	id, err := controller.AddDestination(ctx, dest)
	if err != nil {
		t.Fatalf("AddDestination failed: %v", err)
	}
	t.Logf("Added destination: %s", id)

	got, err := controller.GetDestination(ctx, id)
	if err != nil {
		t.Fatalf("GetDestination failed: %v", err)
	}
	if got.Hostname != dest.Hostname {
		t.Fatalf("Hostname mismatch: got %q, want %q", got.Hostname, dest.Hostname)
	}
	if got.Transport.String() != dest.Transport.String() {
		t.Fatalf("Transport mismatch: got %q, want %q", got.Transport.String(), dest.Transport.String())
	}

	dest.Port = "10514"
	dest.Description = "opnsense-go acceptance test (updated)"
	if err = controller.UpdateDestination(ctx, id, dest); err != nil {
		t.Fatalf("UpdateDestination failed: %v", err)
	}

	got, err = controller.GetDestination(ctx, id)
	if err != nil {
		t.Fatalf("GetDestination (after update) failed: %v", err)
	}
	if got.Port != dest.Port {
		t.Fatalf("Port mismatch after update: got %q, want %q", got.Port, dest.Port)
	}
	if got.Description != dest.Description {
		t.Fatalf("Description mismatch after update: got %q, want %q", got.Description, dest.Description)
	}

	if err = controller.DeleteDestination(ctx, id); err != nil {
		t.Fatalf("DeleteDestination failed: %v", err)
	}
	t.Logf("Deleted destination: %s", id)
}
