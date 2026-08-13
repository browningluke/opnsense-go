package crowdsec

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestGeneral(t *testing.T) {
	opnsense_url := os.Getenv("OPNSENSE_URI")
	opnsense_key := os.Getenv("OPNSENSE_API_KEY")
	opnsense_secret := os.Getenv("OPNSENSE_API_SECRET")

	apiClient := api.NewClient(api.Options{
		Uri:           opnsense_url,
		APIKey:        opnsense_key,
		APISecret:     opnsense_secret,
		AllowInsecure: true,
		MaxBackoff:    30,
		MinBackoff:    1,
		MaxRetries:    4,
	})

	controller := Controller{Api: apiClient}
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
