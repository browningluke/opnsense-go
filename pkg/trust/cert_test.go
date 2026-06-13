package trust

import (
	"context"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestCert(t *testing.T) {
	controller := newController()
	ctx := context.Background()

	// First create a CA to sign with
	ca := &Ca{
		Description: "test-ca-for-cert",
		Action:      api.SelectedMap("internal"),
		KeyType:     api.SelectedMap("2048"),
		Digest:      api.SelectedMap("sha256"),
		Lifetime:    "3650",
		Country:     api.SelectedMap("US"),
		CommonName:  "Test CA For Cert",
	}

	caID, err := controller.AddCa(ctx, ca)
	if err != nil {
		t.Fatalf("AddCa failed: %v", err)
	}
	t.Logf("Created CA: uuid=%s", caID)

	caObj, err := controller.GetCa(ctx, caID)
	if err != nil {
		t.Fatalf("GetCa failed: %v", err)
	}

	cert := &Cert{
		Description:        "test-cert",
		CaRef:              api.SelectedMap(caObj.RefId),
		Action:             api.SelectedMap("internal"),
		KeyType:            api.SelectedMap("2048"),
		Digest:             api.SelectedMap("sha256"),
		CertType:           api.SelectedMap("server_cert"),
		Lifetime:           "397",
		PrivateKeyLocation: api.SelectedMap("firewall"),
		Country:            api.SelectedMap("US"),
		CommonName:         "server.example.com",
		AltnamesDns:        "server.example.com",
	}

	certID, err := controller.AddCert(ctx, cert)
	if err != nil {
		// cleanup CA
		_ = controller.DeleteCa(ctx, caID)
		t.Fatalf("AddCert failed: %v", err)
	}
	t.Logf("AddCert: uuid=%s", certID)

	gotCert, err := controller.GetCert(ctx, certID)
	if err != nil {
		_ = controller.DeleteCa(ctx, caID)
		t.Fatalf("GetCert failed: %v", err)
	}
	t.Logf("GetCert: %+v", gotCert)
	if gotCert.CommonName != "server.example.com" {
		t.Fatalf("CommonName mismatch: got %s", gotCert.CommonName)
	}
	if gotCert.RefId == "" {
		t.Fatal("RefId should be populated")
	}
	if gotCert.CrtPayload == "" {
		t.Fatal("CrtPayload should be populated")
	}

	cert.Description = "test-cert-updated"
	err = controller.UpdateCert(ctx, certID, cert)
	if err != nil {
		_ = controller.DeleteCert(ctx, certID)
		_ = controller.DeleteCa(ctx, caID)
		t.Fatalf("UpdateCert failed: %v", err)
	}

	gotCert, err = controller.GetCert(ctx, certID)
	if err != nil {
		t.Fatalf("GetCert after update failed: %v", err)
	}
	if gotCert.Description != "test-cert-updated" {
		t.Fatalf("Description not updated: got %s", gotCert.Description)
	}

	err = controller.DeleteCert(ctx, certID)
	if err != nil {
		_ = controller.DeleteCa(ctx, caID)
		t.Fatalf("DeleteCert failed: %v", err)
	}
	t.Log("DeleteCert: deleted")

	err = controller.DeleteCa(ctx, caID)
	if err != nil {
		t.Fatalf("DeleteCa failed: %v", err)
	}
	t.Log("DeleteCa: deleted")
}
