package openvpn

import (
	"context"
	"testing"
)

func TestInstance(t *testing.T) {
	c := newTestController(t)
	ctx := context.Background()

	resource := &Instance{
		Enabled:              "1",
		Role:                 "server",
		VPNID:                "8001",
		Description:          "acctest-instance",
		DevType:              "tun",
		Protocol:             "udp",
		Topology:             "subnet",
		VerifyClientCert:     "none",
		AuthMode:             []string{"Local Database"},
		Server:               "10.99.98.0/24",
		Port:                 "11940",
		UsernameAsCommonName: "1",
		RemoteCertTLS:        "0",
		UseOCSP:              "0",
		StrictUserCN:         "0",
		ProvisionExclusive:   "0",
		RegisterDNS:          "0",
	}

	id, err := c.AddInstance(ctx, resource)
	if err != nil {
		t.Fatalf("AddInstance failed: %v", err)
	}
	t.Cleanup(func() {
		if err := c.DeleteInstance(ctx, id); err != nil {
			t.Logf("cleanup DeleteInstance: %v", err)
		}
	})

	got, err := c.GetInstance(ctx, id)
	if err != nil {
		t.Fatalf("GetInstance failed: %v", err)
	}
	if got.Description != "acctest-instance" {
		t.Fatalf("description: got %q, want %q", got.Description, "acctest-instance")
	}
	if got.Role.String() != "server" {
		t.Fatalf("role: got %q, want %q", got.Role.String(), "server")
	}
	if got.Server != "10.99.98.0/24" {
		t.Fatalf("server: got %q, want %q", got.Server, "10.99.98.0/24")
	}
	if got.Port != "11940" {
		t.Fatalf("port: got %q, want %q", got.Port, "11940")
	}

	got.Description = "acctest-instance-upd"
	got.Port = "11941"
	if err := c.UpdateInstance(ctx, id, got); err != nil {
		t.Fatalf("UpdateInstance failed: %v", err)
	}

	got2, err := c.GetInstance(ctx, id)
	if err != nil {
		t.Fatalf("GetInstance after update failed: %v", err)
	}
	if got2.Description != "acctest-instance-upd" {
		t.Fatalf("description after update: got %q, want %q", got2.Description, "acctest-instance-upd")
	}
	if got2.Port != "11941" {
		t.Fatalf("port after update: got %q, want %q", got2.Port, "11941")
	}
}
