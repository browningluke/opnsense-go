package ipsec

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestIPsecAuth(t *testing.T) {
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

	ipsec_auth_local := &IPsecAuthLocal{
		Enabled:        "1",
		Connection:     connection,
		Round:          "1",
		Authentication: "psk",
		Id:             "test-auth-local",
		EAPId:          "",
		Certificates:   api.SelectedMapList{},
		PublicKeys:     api.SelectedMapList{},
		Description:    "Test IPsec Auth Local",
	}

	local_auth_key, err := controller.AddIPsecAuthLocal(ctx, ipsec_auth_local)
	if err != nil {
		t.Fatalf("Failed to add IPsec Auth Local: %v", err)
	}
	t.Logf("Added IPsec Auth Local with ID: %s", local_auth_key)

	retrieved_auth, err := controller.GetIPsecAuthLocal(ctx, local_auth_key)
	if err != nil {
		t.Fatalf("Failed to get IPsec Auth Local: %v", err)
	}
	t.Logf("Retrieved IPsec Auth Local: %+v", retrieved_auth)

	// check all values to see if the retrieved auth matches the original
	if retrieved_auth.Id != ipsec_auth_local.Id {
		t.Errorf("Retrieved ID %s does not match original ID %s", retrieved_auth.Id, ipsec_auth_local.Id)
	}
	if retrieved_auth.Description != ipsec_auth_local.Description {
		t.Errorf("Retrieved Description %s does not match original Description %s", retrieved_auth.Description, ipsec_auth_local.Description)
	}
	if retrieved_auth.Enabled != ipsec_auth_local.Enabled {
		t.Errorf("Retrieved Enabled %s does not match original Enabled %s", retrieved_auth.Enabled, ipsec_auth_local.Enabled)
	}
	if retrieved_auth.Round != ipsec_auth_local.Round {
		t.Errorf("Retrieved Round %s does not match original Round %s", retrieved_auth.Round, ipsec_auth_local.Round)
	}
	if retrieved_auth.Authentication != ipsec_auth_local.Authentication {
		t.Errorf("Retrieved Authentication %s does not match original Authentication %s", retrieved_auth.Authentication, ipsec_auth_local.Authentication)
	}
	if retrieved_auth.Connection != ipsec_auth_local.Connection {
		t.Errorf("Retrieved Connection %s does not match original Connection %s", retrieved_auth.Connection, ipsec_auth_local.Connection)
	}
	if !sliceEqual(retrieved_auth.Certificates, ipsec_auth_local.Certificates) {
		t.Errorf("Retrieved Certificates %v does not match original Certificates %v", retrieved_auth.Certificates, ipsec_auth_local.Certificates)
	}
	if !sliceEqual(retrieved_auth.PublicKeys, ipsec_auth_local.PublicKeys) {
		t.Errorf("Retrieved PublicKeys %v does not match original PublicKeys %v", retrieved_auth.PublicKeys, ipsec_auth_local.PublicKeys)
	}
	if retrieved_auth.EAPId != ipsec_auth_local.EAPId {
		t.Errorf("Retrieved EAPId %s does not match original EAPId %s", retrieved_auth.EAPId, ipsec_auth_local.EAPId)
	}

	// Update the IPsec Auth Local
	ipsec_auth_local.Description = "Updated Test IPsec Auth Local"
	err = controller.UpdateIPsecAuthLocal(ctx, local_auth_key, ipsec_auth_local)
	if err != nil {
		t.Fatalf("Failed to update IPsec Auth Local: %v", err)
	}
	t.Logf("Updated IPsec Auth Local with ID: %s", local_auth_key)

	err = controller.DeleteIPsecAuthLocal(ctx, local_auth_key)
	if err != nil {
		t.Fatalf("Failed to delete IPsec Auth Local: %v", err)
	}

	err = controller.DeleteIPsecConnection(ctx, conn_key)
	if err != nil {
		t.Fatalf("Failed to delete IPsec Connection: %v", err)
	}
	t.Logf("Deleted IPsec Connection with key: %s", conn_key)
}
