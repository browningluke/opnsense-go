package wireguard

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
	if resp.General.Enabled == "" {
		t.Fatal("expected general.enabled to be set")
	}

	origEnabled := resp.General.Enabled

	// Step 2: Toggle enabled (flip its value).
	newEnabled := "1"
	if origEnabled == "1" {
		newEnabled = "0"
	}
	setResp, err := controller.GeneralSet(ctx, &WireguardGeneral{Enabled: newEnabled})
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
	if resp2.General.Enabled != newEnabled {
		t.Fatalf("enabled not updated; got %q, want %q", resp2.General.Enabled, newEnabled)
	}

	// Step 4: Restore original value.
	_, err = controller.GeneralSet(ctx, &WireguardGeneral{Enabled: origEnabled})
	if err != nil {
		t.Fatalf("GeneralSet (restore) failed: %v", err)
	}
}

func TestServiceGenKeyPair(t *testing.T) {
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

	resp, err := controller.ServiceGenKeyPair(ctx)
	if err != nil {
		t.Fatalf("ServiceGenKeyPair failed: %v", err)
	}
	t.Logf("ServiceGenKeyPair: %+v", resp)

	if resp.PrivKey == "" {
		t.Fatal("expected privkey to be set")
	}
	if resp.PubKey == "" {
		t.Fatal("expected pubkey to be set")
	}
	if resp.Status != "ok" {
		t.Fatalf("expected status=ok, got %q", resp.Status)
	}

	// Verify two consecutive calls return different keys.
	resp2, err := controller.ServiceGenKeyPair(ctx)
	if err != nil {
		t.Fatalf("ServiceGenKeyPair (second call) failed: %v", err)
	}
	if resp.PrivKey == resp2.PrivKey {
		t.Fatal("expected consecutive calls to return different private keys")
	}
}

func TestServiceGenPsk(t *testing.T) {
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

	resp, err := controller.ServiceGenPsk(ctx)
	if err != nil {
		t.Fatalf("ServiceGenPsk failed: %v", err)
	}
	t.Logf("ServiceGenPsk: %+v", resp)

	if resp.Psk == "" {
		t.Fatal("expected psk to be set")
	}
	if resp.Status != "ok" {
		t.Fatalf("expected status=ok, got %q", resp.Status)
	}

	// Verify two consecutive calls return different PSKs.
	resp2, err := controller.ServiceGenPsk(ctx)
	if err != nil {
		t.Fatalf("ServiceGenPsk (second call) failed: %v", err)
	}
	if resp.Psk == resp2.Psk {
		t.Fatal("expected consecutive calls to return different PSKs")
	}
}
