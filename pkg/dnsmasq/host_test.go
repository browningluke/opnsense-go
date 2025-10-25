package dnsmasq

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestHost(t *testing.T) {
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

	host := &Host{
		Hostname:          "test-host",
		IpAddresses:       api.SelectedMapList([]string{"192.168.2.50"}),
		AliasRecords:      api.SelectedMapList([]string{"test-alias"}),
		CnameRecords:      api.SelectedMapList([]string{"test-cname"}),
		HardwareAddresses: api.SelectedMapList([]string{"AA:BB:CC:DD:EE:FF"}),
		// Tagset:            api.SelectedMap("4a7a61e0-9be1-49e5-86bb-7f11ef274764"),
	}

	respAdd, err := controller.AddHost(ctx, host)
	if err != nil {
		t.Fatalf("Failed to add host: %v", err)
	}
	t.Logf("AddHost: %+v", respAdd)

	respGet, err := controller.GetHost(ctx, respAdd)
	if err != nil {
		t.Fatalf("Failed to get host: %v", err)
	}
	t.Logf("GetHost: %+v", respGet)

	host.Hostname = "test-host-updated"
	host.IpAddresses = api.SelectedMapList([]string{"192.168.2.51"})
	host.AliasRecords = api.SelectedMapList([]string{"test-alias-updated"})
	host.CnameRecords = api.SelectedMapList([]string{"test-cname-updated"})
	host.HardwareAddresses = api.SelectedMapList([]string{"00:11:22:33:44:55"})
	// host.Tagset = api.SelectedMap("ee791e64-69de-4776-a9a0-e1442630c6ef")
	err = controller.UpdateHost(ctx, respAdd, host)
	if err != nil {
		t.Fatalf("Failed to update host: %v", err)
	}
	t.Logf("UpdateHost: %+v", host)

	respGet, err = controller.GetHost(ctx, respAdd)
	if err != nil {
		t.Fatalf("Failed to get host: %v", err)
	}
	t.Logf("GetHost: %+v", respGet)

	err = controller.DeleteHost(ctx, respAdd)
	if err != nil {
		t.Fatalf("Failed to delete host: %v", err)
	}
	t.Log("DeleteHost: Deleted!")
}
