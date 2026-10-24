package mdnsrepeater

import (
	"context"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestSettings(t *testing.T) {
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

	resp, err := controller.SettingsGet(ctx)
	if err != nil && strings.Contains(err.Error(), "status code 404") {
		t.Skip("os-mdns-repeater plugin is not installed")
	}
	if err != nil {
		t.Fatalf("SettingsGet failed: %v", err)
	}
	orig := resp.MDNSRepeater
	t.Cleanup(func() {
		if _, err := controller.SettingsSet(ctx, &orig); err != nil {
			t.Errorf("SettingsSet (restore) failed: %v", err)
			return
		}
		// The test stops the service; restart it if it was running.
		if _, err := controller.SettingsReconfigure(ctx); err != nil {
			t.Errorf("SettingsReconfigure (restore) failed: %v", err)
		}
	})

	// The CI VM only has a wan interface. Leave the repeater disabled: it
	// refuses to start on fewer than two interfaces.
	want := Settings{
		Enabled:    "0",
		EnableCarp: "1",
		Interfaces: api.SelectedMapList{"wan"},
		Blocklist:  api.SelectedMapList{"10.0.0.0/8", "192.168.1.0/24"},
	}
	setResp, err := controller.SettingsSet(ctx, &want)
	if err != nil {
		t.Fatalf("SettingsSet failed: %v", err)
	}
	if setResp.Result != "saved" {
		t.Fatalf("expected result=saved, got %+v", setResp)
	}

	rcResp, err := controller.SettingsReconfigure(ctx)
	if err != nil {
		t.Fatalf("SettingsReconfigure failed: %v", err)
	}
	if rcResp.Status != "ok" {
		t.Fatalf("expected status=ok, got %q", rcResp.Status)
	}

	resp, err = controller.SettingsGet(ctx)
	if err != nil {
		t.Fatalf("SettingsGet after update failed: %v", err)
	}
	got := resp.MDNSRepeater
	slices.Sort(got.Blocklist)
	if got.Enabled != want.Enabled ||
		got.EnableCarp != want.EnableCarp ||
		!slices.Equal(got.Interfaces, want.Interfaces) ||
		!slices.Equal(got.Blocklist, want.Blocklist) {
		t.Fatalf("settings not updated; got %+v, want %+v", got, want)
	}
}
