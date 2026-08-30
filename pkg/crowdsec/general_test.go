package crowdsec

import (
	"context"
	"testing"
)

func TestGeneral(t *testing.T) {
	controller := newController()
	ctx := context.Background()

	// Step 1: Get current settings.
	resp, err := controller.GeneralGet(ctx)
	if err != nil {
		t.Fatalf("GeneralGet failed: %v", err)
	}
	t.Logf("GeneralGet: %+v", resp)
	if resp.General.LapiListenAddress == "" {
		t.Fatal("expected general.lapi_listen_address to be set")
	}

	orig := resp.General

	// Step 2: Toggle rules_log (flip its value) and write back.
	newRulesLog := "1"
	if orig.RulesLog == "1" {
		newRulesLog = "0"
	}
	updated := orig
	updated.RulesLog = newRulesLog

	setResp, err := controller.GeneralSet(ctx, &updated)
	if err != nil {
		t.Fatalf("GeneralSet failed: %v", err)
	}
	t.Logf("GeneralSet (toggle): %+v", setResp)
	if setResp.Result != "saved" {
		t.Fatalf("expected result=saved, got %q", setResp.Result)
	}

	// Step 3: Read back and verify.
	resp2, err := controller.GeneralGet(ctx)
	if err != nil {
		t.Fatalf("GeneralGet after update failed: %v", err)
	}
	if resp2.General.RulesLog != newRulesLog {
		t.Fatalf("rules_log not updated; got %q, want %q", resp2.General.RulesLog, newRulesLog)
	}

	// Step 4: Restore original value.
	_, err = controller.GeneralSet(ctx, &orig)
	if err != nil {
		t.Fatalf("GeneralSet (restore) failed: %v", err)
	}
}
