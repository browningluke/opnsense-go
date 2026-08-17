package monit

import (
	"context"
	"testing"
)

func TestSettings(t *testing.T) {
	controller := newController()
	ctx := context.Background()

	result, err := controller.SettingsGet(ctx)
	if err != nil {
		t.Fatalf("SettingsGet failed: %v", err)
	}
	t.Logf("SettingsGet: %+v", result.Monit)

	// Update settings (no-op, set same values back)
	_, err = controller.SettingsSet(ctx, &result.Monit)
	if err != nil {
		t.Fatalf("SettingsSet failed: %v", err)
	}

	_, err = controller.SettingsReconfigure(ctx)
	if err != nil {
		t.Fatalf("SettingsReconfigure failed: %v", err)
	}
	t.Log("SettingsReconfigure: ok")
}
