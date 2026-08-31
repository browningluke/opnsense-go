package openvpn

import (
	"context"
	"testing"
)

func TestClientOverwrite(t *testing.T) {
	c := newTestController(t)
	ctx := context.Background()

	resource := &ClientOverwrite{
		Enabled:     "1",
		CommonName:  "acctest-client",
		Description: "acctest-cso",
		Block:       "0",
		PushReset:   "0",
		RegisterDNS: "0",
	}

	id, err := c.AddClientOverwrite(ctx, resource)
	if err != nil {
		t.Fatalf("AddClientOverwrite failed: %v", err)
	}
	t.Cleanup(func() {
		if err := c.DeleteClientOverwrite(ctx, id); err != nil {
			t.Logf("cleanup DeleteClientOverwrite: %v", err)
		}
	})

	got, err := c.GetClientOverwrite(ctx, id)
	if err != nil {
		t.Fatalf("GetClientOverwrite failed: %v", err)
	}
	if got.CommonName != "acctest-client" {
		t.Fatalf("common_name: got %q", got.CommonName)
	}
	if got.Description != "acctest-cso" {
		t.Fatalf("description: got %q", got.Description)
	}

	got.Description = "acctest-cso-upd"
	got.TunnelNetwork = "10.50.0.0/24"
	if err := c.UpdateClientOverwrite(ctx, id, got); err != nil {
		t.Fatalf("UpdateClientOverwrite failed: %v", err)
	}

	got2, err := c.GetClientOverwrite(ctx, id)
	if err != nil {
		t.Fatalf("GetClientOverwrite after update failed: %v", err)
	}
	if got2.Description != "acctest-cso-upd" {
		t.Fatalf("description after update: got %q", got2.Description)
	}
	if got2.TunnelNetwork != "10.50.0.0/24" {
		t.Fatalf("tunnel_network after update: got %q", got2.TunnelNetwork)
	}
}
