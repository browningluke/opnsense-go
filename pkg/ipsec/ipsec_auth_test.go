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

	retrieved_auth_local, err := controller.GetIPsecAuthLocal(ctx, local_auth_key)
	if err != nil {
		t.Fatalf("Failed to get IPsec Auth Local: %v", err)
	}
	t.Logf("Retrieved IPsec Auth Local: %+v", retrieved_auth_local)

	if retrieved_auth_local.Id != ipsec_auth_local.Id {
		t.Errorf("Retrieved ID %s does not match original ID %s", retrieved_auth_local.Id, ipsec_auth_local.Id)
	}
	if retrieved_auth_local.Description != ipsec_auth_local.Description {
		t.Errorf("Retrieved Description %s does not match original Description %s", retrieved_auth_local.Description, ipsec_auth_local.Description)
	}
	if retrieved_auth_local.Enabled != ipsec_auth_local.Enabled {
		t.Errorf("Retrieved Enabled %s does not match original Enabled %s", retrieved_auth_local.Enabled, ipsec_auth_local.Enabled)
	}
	if retrieved_auth_local.Round != ipsec_auth_local.Round {
		t.Errorf("Retrieved Round %s does not match original Round %s", retrieved_auth_local.Round, ipsec_auth_local.Round)
	}
	if retrieved_auth_local.Authentication != ipsec_auth_local.Authentication {
		t.Errorf("Retrieved Authentication %s does not match original Authentication %s", retrieved_auth_local.Authentication, ipsec_auth_local.Authentication)
	}
	if retrieved_auth_local.Connection != ipsec_auth_local.Connection {
		t.Errorf("Retrieved Connection %s does not match original Connection %s", retrieved_auth_local.Connection, ipsec_auth_local.Connection)
	}
	if !sliceEqual(retrieved_auth_local.Certificates, ipsec_auth_local.Certificates) {
		t.Errorf("Retrieved Certificates %v does not match original Certificates %v", retrieved_auth_local.Certificates, ipsec_auth_local.Certificates)
	}
	if !sliceEqual(retrieved_auth_local.PublicKeys, ipsec_auth_local.PublicKeys) {
		t.Errorf("Retrieved PublicKeys %v does not match original PublicKeys %v", retrieved_auth_local.PublicKeys, ipsec_auth_local.PublicKeys)
	}
	if retrieved_auth_local.EAPId != ipsec_auth_local.EAPId {
		t.Errorf("Retrieved EAPId %s does not match original EAPId %s", retrieved_auth_local.EAPId, ipsec_auth_local.EAPId)
	}

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

	ipsec_auth_remote := &IPsecAuthRemote{
		Enabled:        "1",
		Connection:     connection,
		Round:          "1",
		Authentication: "psk",
		Id:             "test-auth-remote",
		EAPId:          "",
		Certificates:   api.SelectedMapList{},
		PublicKeys:     api.SelectedMapList{},
		Description:    "Test IPsec Auth Remote",
	}

	remote_auth_key, err := controller.AddIPsecAuthRemote(ctx, ipsec_auth_remote)
	if err != nil {
		t.Fatalf("Failed to add IPsec Auth Remote: %v", err)
	}
	t.Logf("Added IPsec Auth Remote with ID: %s", remote_auth_key)

	retrieved_auth_remote, err := controller.GetIPsecAuthRemote(ctx, remote_auth_key)
	if err != nil {
		t.Fatalf("Failed to get IPsec Auth Local: %v", err)
	}
	t.Logf("Retrieved IPsec Auth Local: %+v", retrieved_auth_remote)

	if retrieved_auth_remote.Id != ipsec_auth_remote.Id {
		t.Errorf("Retrieved ID %s does not match original ID %s", retrieved_auth_remote.Id, ipsec_auth_remote.Id)
	}
	if retrieved_auth_remote.Description != ipsec_auth_remote.Description {
		t.Errorf("Retrieved Description %s does not match original Description %s", retrieved_auth_remote.Description, ipsec_auth_remote.Description)
	}
	if retrieved_auth_remote.Enabled != ipsec_auth_remote.Enabled {
		t.Errorf("Retrieved Enabled %s does not match original Enabled %s", retrieved_auth_remote.Enabled, ipsec_auth_remote.Enabled)
	}
	if retrieved_auth_remote.Round != ipsec_auth_remote.Round {
		t.Errorf("Retrieved Round %s does not match original Round %s", retrieved_auth_remote.Round, ipsec_auth_remote.Round)
	}
	if retrieved_auth_remote.Authentication != ipsec_auth_remote.Authentication {
		t.Errorf("Retrieved Authentication %s does not match original Authentication %s", retrieved_auth_remote.Authentication, ipsec_auth_remote.Authentication)
	}
	if retrieved_auth_remote.Connection != ipsec_auth_remote.Connection {
		t.Errorf("Retrieved Connection %s does not match original Connection %s", retrieved_auth_remote.Connection, ipsec_auth_remote.Connection)
	}
	if !sliceEqual(retrieved_auth_remote.Certificates, ipsec_auth_remote.Certificates) {
		t.Errorf("Retrieved Certificates %v does not match original Certificates %v", retrieved_auth_remote.Certificates, ipsec_auth_remote.Certificates)
	}
	if !sliceEqual(retrieved_auth_remote.PublicKeys, ipsec_auth_remote.PublicKeys) {
		t.Errorf("Retrieved PublicKeys %v does not match original PublicKeys %v", retrieved_auth_remote.PublicKeys, ipsec_auth_remote.PublicKeys)
	}
	if retrieved_auth_remote.EAPId != ipsec_auth_remote.EAPId {
		t.Errorf("Retrieved EAPId %s does not match original EAPId %s", retrieved_auth_remote.EAPId, ipsec_auth_remote.EAPId)
	}

	ipsec_auth_remote.Description = "Updated Test IPsec Auth Remote"
	err = controller.UpdateIPsecAuthRemote(ctx, remote_auth_key, ipsec_auth_remote)
	if err != nil {
		t.Fatalf("Failed to update IPsec Auth Remote: %v", err)
	}
	t.Logf("Updated IPsec Auth Remote with ID: %s", remote_auth_key)

	err = controller.DeleteIPsecAuthRemote(ctx, remote_auth_key)
	if err != nil {
		t.Fatalf("Failed to delete IPsec Auth Remote: %v", err)
	}

	err = controller.DeleteIPsecConnection(ctx, conn_key)
	if err != nil {
		t.Fatalf("Failed to delete IPsec Connection: %v", err)
	}
	t.Logf("Deleted IPsec Connection with key: %s", conn_key)
}
