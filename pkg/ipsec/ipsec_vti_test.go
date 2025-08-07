package ipsec

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestIPsecVTI(t *testing.T) {
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

	proposals := api.SelectedMapList{
		"aes256-sha256-modp2048",
	}
	local_addresses := api.SelectedMapList{
		"192.168.1.1",
	}
	remote_addresses := api.SelectedMapList{
		"192.168.1.2",
	}
	ipsec_connection := &IPsecConnection{
		Enabled:                "1",
		Proposals:              proposals,
		Unique:                 "no",
		Aggressive:             "0",
		Version:                "2",
		Mobike:                 "1",
		LocalAddresses:         local_addresses,
		RemoteAddresses:        remote_addresses,
		LocalPort:              "",
		RemotePort:             "",
		UDPEncapsulation:       "0",
		ReauthenticationTime:   "0",
		RekeyTime:              "28800",
		IKELifetime:            "86400",
		DPDDelay:               "120",
		DPDTimeout:             "540",
		IPPools:                api.SelectedMapList{},
		SendCertificateRequest: "1",
		SendCertificate:        "always",
		KeyingTries:            "3",
		Description:            "Test IPsec Connection",
	}
	conn_key, err := controller.AddIPsecConnection(ctx, ipsec_connection)
	if err != nil {
		t.Fatalf("Failed to add IPsec Connection: %v", err)
	}
	t.Logf("Added IPsec Connection with key: %s", conn_key)

	ipsec_vti := &IPsecVTI{
		Enabled:         "1",
		RequestID:       "1234",
		LocalIP:         "2.3.4.5",
		RemoteIP:        "5.6.7.8",
		TunnelLocalIP:   "1.2.3.4",
		TunnelRemoteIP:  "4.3.2.1",
		TunnelLocalIP2:  "7.8.9.10",
		TunnelRemoteIP2: "8.7.6.5",
		Description:     "Test IPsec VTI",
	}

	vti_key, err := controller.AddIPsecVTI(ctx, ipsec_vti)
	if err != nil {
		t.Fatalf("Failed to add IPsec VTI: %v", err)
	}
	t.Logf("Added IPsec VTI with key: %s", vti_key)

	retrieved_vti, err := controller.GetIPsecVTI(ctx, vti_key)
	if err != nil {
		t.Fatalf("Failed to get IPsec VTI: %v", err)
	}
	if retrieved_vti.Description != ipsec_vti.Description {
		t.Errorf("Expected description '%s', got '%s'", ipsec_vti.Description, retrieved_vti.Description)
	}
	if retrieved_vti.LocalIP != ipsec_vti.LocalIP {
		t.Errorf("Expected LocalIP '%s', got '%s'", ipsec_vti.LocalIP, retrieved_vti.LocalIP)
	}
	if retrieved_vti.RemoteIP != ipsec_vti.RemoteIP {
		t.Errorf("Expected RemoteIP '%s', got '%s'", ipsec_vti.RemoteIP, retrieved_vti.RemoteIP)
	}
	if retrieved_vti.TunnelLocalIP != ipsec_vti.TunnelLocalIP {
		t.Errorf("Expected TunnelLocalIP '%s', got '%s'", ipsec_vti.TunnelLocalIP, retrieved_vti.TunnelLocalIP)
	}
	if retrieved_vti.TunnelRemoteIP != ipsec_vti.TunnelRemoteIP {
		t.Errorf("Expected TunnelRemoteIP '%s', got '%s'", ipsec_vti.TunnelRemoteIP, retrieved_vti.TunnelRemoteIP)
	}
	if retrieved_vti.TunnelLocalIP2 != ipsec_vti.TunnelLocalIP2 {
		t.Errorf("Expected TunnelLocalIP2 '%s', got '%s'", ipsec_vti.TunnelLocalIP2, retrieved_vti.TunnelLocalIP2)
	}
	if retrieved_vti.TunnelRemoteIP2 != ipsec_vti.TunnelRemoteIP2 {
		t.Errorf("Expected TunnelRemoteIP2 '%s', got '%s'", ipsec_vti.TunnelRemoteIP2, retrieved_vti.TunnelRemoteIP2)
	}
	if retrieved_vti.RequestID != ipsec_vti.RequestID {
		t.Errorf("Expected RequestID '%s', got '%s'", ipsec_vti.RequestID, retrieved_vti.RequestID)
	}
	if retrieved_vti.Enabled != ipsec_vti.Enabled {
		t.Errorf("Expected Enabled '%s', got '%s'", ipsec_vti.Enabled, retrieved_vti.Enabled)
	}
	ipsec_vti.Description = "Updated Test IPsec VTI"
	t.Logf("Updating IPsec VTI with key: %s", vti_key)
	err = controller.UpdateIPsecVTI(ctx, vti_key, ipsec_vti)
	if err != nil {
		t.Fatalf("Failed to update IPsec VTI: %v", err)
	}
	t.Logf("Updated IPsec VTI with key: %s", vti_key)
	err = controller.DeleteIPsecVTI(ctx, vti_key)
	if err != nil {
		t.Fatalf("Failed to delete IPsec VTI: %v", err)
	}
	t.Logf("Deleted IPsec VTI with key: %s", vti_key)

	err = controller.DeleteIPsecConnection(ctx, conn_key)
	if err != nil {
		t.Fatalf("Failed to delete IPsec Connection: %v", err)
	}
	t.Logf("Deleted IPsec Connection with key: %s", conn_key)
}
