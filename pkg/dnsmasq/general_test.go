package dnsmasq

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

	respGet, err := controller.GeneralSettingsGet(ctx)
	if err != nil {
		t.Fatalf("Failed to get general settings: %v", err)
	}
	t.Logf("GeneralSettingsGet: %+v", respGet)

	// Default settings from out-of-box install opnsense 25.7
	reqObj := &GeneralSettings{
		IsEnabled:   "1",
		DNS_Port:    "0",
		DNS_NoIdent: "1",
		DHCPSettings: GeneralDHCPSettings{
			FQDN:                  "1",
			RegisterFirewallRules: "1",
			DisableHASync:         "1", // New setting to set without breaking anything
		},
	}

	respSet, err := controller.GeneralSettingsSet(ctx, reqObj)
	if err != nil {
		t.Fatalf("Failed to set general settings: %v", err)
	}
	t.Logf("GeneralSettingsSet: %+v", respSet)

	respGet, err = controller.GeneralSettingsGet(ctx)
	if err != nil {
		t.Fatalf("Failed to get general settings: %v", err)
	}
	t.Logf("GeneralSettingsGet: %+v", respGet)

	if respGet.Dnsmasq.DHCPSettings.DisableHASync != reqObj.DHCPSettings.DisableHASync {
		t.Fatalf("Failed to set DisableHASync")
	}

	// Resetting
	// Default settings from out-of-box install opnsense 25.7
	reqObj = &GeneralSettings{
		IsEnabled:   "1",
		DNS_Port:    "0",
		DNS_NoIdent: "1",
		DHCPSettings: GeneralDHCPSettings{
			FQDN:                  "1",
			RegisterFirewallRules: "1",
			DisableHASync:         "0", // New setting to set without breaking anything
		},
	}

	respSet, err = controller.GeneralSettingsSet(ctx, reqObj)
	if err != nil {
		t.Fatalf("Failed to set general settings: %v", err)
	}
	t.Logf("GeneralSettingsSet: %+v", respSet)
}
