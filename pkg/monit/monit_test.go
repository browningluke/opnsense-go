package monit

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func newController() Controller {
	return Controller{
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

func TestTest(t *testing.T) {
	c := newController()
	ctx := context.Background()

	tst := &Test{
		Name:      "opnsense-go-acc-test",
		Type:      api.SelectedMap("ProgramStatus"),
		Condition: "status != 0",
		Action:    api.SelectedMap("alert"),
		Path:      "",
	}

	id, err := c.AddTest(ctx, tst)
	if err != nil {
		t.Fatalf("AddTest failed: %v", err)
	}
	t.Logf("Added test: %s", id)

	got, err := c.GetTest(ctx, id)
	if err != nil {
		t.Fatalf("GetTest failed: %v", err)
	}
	if got.Name != tst.Name {
		t.Fatalf("Name mismatch: got %q, want %q", got.Name, tst.Name)
	}
	if got.Type.String() != tst.Type.String() {
		t.Fatalf("Type mismatch: got %q, want %q", got.Type.String(), tst.Type.String())
	}

	tst.Condition = "status != 1"
	if err = c.UpdateTest(ctx, id, tst); err != nil {
		t.Fatalf("UpdateTest failed: %v", err)
	}

	got, err = c.GetTest(ctx, id)
	if err != nil {
		t.Fatalf("GetTest (after update) failed: %v", err)
	}
	if got.Condition != tst.Condition {
		t.Fatalf("Condition mismatch after update: got %q, want %q", got.Condition, tst.Condition)
	}

	if err = c.DeleteTest(ctx, id); err != nil {
		t.Fatalf("DeleteTest failed: %v", err)
	}
	t.Logf("Deleted test: %s", id)
}

func TestService(t *testing.T) {
	c := newController()
	ctx := context.Background()

	svc := &Service{
		Enabled:      "1",
		Name:         "opnsense-go-acc-service",
		Description:  "opnsense-go acceptance test",
		Type:         api.SelectedMap("process"),
		Pidfile:      "/var/run/sshd.pid",
		Match:        "",
		Path:         "",
		Timeout:      "300",
		Starttimeout: "30",
		Address:      "",
		Interface:    api.SelectedMap(""),
		Start:        "",
		Stop:         "",
		Tests:        api.SelectedMapList{},
		Depends:      api.SelectedMapList{},
		Polltime:     "",
	}

	id, err := c.AddService(ctx, svc)
	if err != nil {
		t.Fatalf("AddService failed: %v", err)
	}
	t.Logf("Added service: %s", id)

	got, err := c.GetService(ctx, id)
	if err != nil {
		t.Fatalf("GetService failed: %v", err)
	}
	if got.Name != svc.Name {
		t.Fatalf("Name mismatch: got %q, want %q", got.Name, svc.Name)
	}
	if got.Type.String() != svc.Type.String() {
		t.Fatalf("Type mismatch: got %q, want %q", got.Type.String(), svc.Type.String())
	}

	svc.Description = "opnsense-go acceptance test (updated)"
	svc.Timeout = "60"
	if err = c.UpdateService(ctx, id, svc); err != nil {
		t.Fatalf("UpdateService failed: %v", err)
	}

	got, err = c.GetService(ctx, id)
	if err != nil {
		t.Fatalf("GetService (after update) failed: %v", err)
	}
	if got.Timeout != svc.Timeout {
		t.Fatalf("Timeout mismatch after update: got %q, want %q", got.Timeout, svc.Timeout)
	}

	if err = c.DeleteService(ctx, id); err != nil {
		t.Fatalf("DeleteService failed: %v", err)
	}
	t.Logf("Deleted service: %s", id)
}

func TestAlert(t *testing.T) {
	c := newController()
	ctx := context.Background()

	alert := &Alert{
		Enabled:     "1",
		Recipient:   "admin@example.com",
		Noton:       "0",
		Events:      api.SelectedMapList{},
		Format:      "",
		Reminder:    "",
		Description: "opnsense-go acceptance test",
	}

	id, err := c.AddAlert(ctx, alert)
	if err != nil {
		t.Fatalf("AddAlert failed: %v", err)
	}
	t.Logf("Added alert: %s", id)

	got, err := c.GetAlert(ctx, id)
	if err != nil {
		t.Fatalf("GetAlert failed: %v", err)
	}
	if got.Recipient != alert.Recipient {
		t.Fatalf("Recipient mismatch: got %q, want %q", got.Recipient, alert.Recipient)
	}

	alert.Description = "opnsense-go acceptance test (updated)"
	if err = c.UpdateAlert(ctx, id, alert); err != nil {
		t.Fatalf("UpdateAlert failed: %v", err)
	}

	got, err = c.GetAlert(ctx, id)
	if err != nil {
		t.Fatalf("GetAlert (after update) failed: %v", err)
	}
	if got.Description != alert.Description {
		t.Fatalf("Description mismatch after update: got %q, want %q", got.Description, alert.Description)
	}

	if err = c.DeleteAlert(ctx, id); err != nil {
		t.Fatalf("DeleteAlert failed: %v", err)
	}
	t.Logf("Deleted alert: %s", id)
}
