package ipsec

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestIPsecChild(t *testing.T) {
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

	connection := api.SelectedMap(conn_key)

	ipsec_child := &IPsecChild{
		Enabled:         "1",
		Connection:      connection,
		Proposals:       proposals,
		SHA256_96:       "0",
		StartAction:     api.SelectedMap("none"),
		CloseAction:     api.SelectedMap("none"),
		DPDAction:       api.SelectedMap("clear"),
		Mode:            api.SelectedMap("tunnel"),
		InstallPolicies: "1",
		LocalNetworks:   local_addresses,
		RemoteNetworks:  remote_addresses,
		RequestID:       "1",
		RekeyTime:       "3600",
		Description:     "Test IPsec Child",
	}

	child_key, err := controller.AddIPsecChild(ctx, ipsec_child)
	if err != nil {
		t.Fatalf("Failed to add IPsec Child: %v", err)
	}
	t.Logf("Added IPsec Child with key: %s", child_key)

	retrieved_child, err := controller.GetIPsecChild(ctx, child_key)
	if err != nil {
		t.Fatalf("Failed to get IPsec Child: %v", err)
	}
	t.Logf("Retrieved IPsec Child: %+v", retrieved_child)

	if retrieved_child.Description != ipsec_child.Description {
		t.Errorf("Expected description '%s', got '%s'", ipsec_child.Description, retrieved_child.Description)
	}
	if retrieved_child.Enabled != ipsec_child.Enabled {
		t.Errorf("Expected enabled '%s', got '%s'", ipsec_child.Enabled, retrieved_child.Enabled)
	}
	if retrieved_child.Connection != ipsec_child.Connection {
		t.Errorf("Expected connection '%s', got '%s'", ipsec_child.Connection, retrieved_child.Connection)
	}
	if !sliceEqual(retrieved_child.Proposals, ipsec_child.Proposals) {
		t.Errorf("Expected proposals '%s', got '%s'", ipsec_child.Proposals, retrieved_child.Proposals)
	}
	if retrieved_child.SHA256_96 != ipsec_child.SHA256_96 {
		t.Errorf("Expected SHA256_96 '%s', got '%s'", ipsec_child.SHA256_96, retrieved_child.SHA256_96)
	}
	if retrieved_child.StartAction != ipsec_child.StartAction {
		t.Errorf("Expected StartAction '%s', got '%s'", ipsec_child.StartAction, retrieved_child.StartAction)
	}
	if retrieved_child.CloseAction != ipsec_child.CloseAction {
		t.Errorf("Expected CloseAction '%s', got '%s'", ipsec_child.CloseAction, retrieved_child.CloseAction)
	}
	if retrieved_child.DPDAction != ipsec_child.DPDAction {
		t.Errorf("Expected DPDAction '%s', got '%s'", ipsec_child.DPDAction, retrieved_child.DPDAction)
	}
	if retrieved_child.Mode != ipsec_child.Mode {
		t.Errorf("Expected Mode '%s', got '%s'", ipsec_child.Mode, retrieved_child.Mode)
	}
	if retrieved_child.InstallPolicies != ipsec_child.InstallPolicies {
		t.Errorf("Expected InstallPolicies '%s', got '%s'", ipsec_child.InstallPolicies, retrieved_child.InstallPolicies)
	}
	if !sliceEqual(retrieved_child.LocalNetworks, ipsec_child.LocalNetworks) {
		t.Errorf("Expected LocalNetworks '%s', got '%s'", ipsec_child.LocalNetworks, retrieved_child.LocalNetworks)
	}
	if !sliceEqual(retrieved_child.RemoteNetworks, ipsec_child.RemoteNetworks) {
		t.Errorf("Expected RemoteNetworks '%s', got '%s'", ipsec_child.RemoteNetworks, retrieved_child.RemoteNetworks)
	}
	if retrieved_child.RequestID != ipsec_child.RequestID {
		t.Errorf("Expected RequestID '%s', got '%s'", ipsec_child.RequestID, retrieved_child.RequestID)
	}
	if retrieved_child.RekeyTime != ipsec_child.RekeyTime {
		t.Errorf("Expected RekeyTime '%s', got '%s'", ipsec_child.RekeyTime, retrieved_child.RekeyTime)
	}
	ipsec_child.Description = "Updated Test IPsec Child"
	err = controller.UpdateIPsecChild(ctx, child_key, ipsec_child)
	if err != nil {
		t.Fatalf("Failed to update IPsec Child: %v", err)
	}
	t.Logf("Updated IPsec Child with key: %s", child_key)

	err = controller.DeleteIPsecChild(ctx, child_key)
	if err != nil {
		t.Fatalf("Failed to delete IPsec Child: %v", err)
	}
	t.Logf("Deleted IPsec Child with key: %s", child_key)

	err = controller.DeleteIPsecConnection(ctx, conn_key)
	if err != nil {
		t.Fatalf("Failed to delete IPsec Connection: %v", err)
	}
	t.Logf("Deleted IPsec Connection with key: %s", conn_key)
}
