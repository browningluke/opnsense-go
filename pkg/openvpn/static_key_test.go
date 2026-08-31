package openvpn

import (
	"context"
	"testing"
)

func TestStaticKey(t *testing.T) {
	c := newTestController(t)
	ctx := context.Background()

	gen, err := c.ServiceGenKey(ctx, nil)
	if err != nil {
		t.Fatalf("ServiceGenKey failed: %v", err)
	}
	if gen.Key == "" {
		t.Fatal("ServiceGenKey returned empty key")
	}

	resource := &StaticKey{
		Mode:        "crypt",
		Key:         gen.Key,
		Description: "acctest-static-key",
	}

	id, err := c.AddStaticKey(ctx, resource)
	if err != nil {
		t.Fatalf("AddStaticKey failed: %v", err)
	}
	t.Cleanup(func() {
		if err := c.DeleteStaticKey(ctx, id); err != nil {
			t.Logf("cleanup DeleteStaticKey: %v", err)
		}
	})

	got, err := c.GetStaticKey(ctx, id)
	if err != nil {
		t.Fatalf("GetStaticKey failed: %v", err)
	}
	if got.Description != "acctest-static-key" {
		t.Fatalf("description: got %q, want %q", got.Description, "acctest-static-key")
	}
	if got.Mode.String() != "crypt" {
		t.Fatalf("mode: got %q, want %q", got.Mode.String(), "crypt")
	}

	got.Description = "acctest-static-key-upd"
	got.Mode = "auth"
	if err := c.UpdateStaticKey(ctx, id, got); err != nil {
		t.Fatalf("UpdateStaticKey failed: %v", err)
	}

	got2, err := c.GetStaticKey(ctx, id)
	if err != nil {
		t.Fatalf("GetStaticKey after update failed: %v", err)
	}
	if got2.Description != "acctest-static-key-upd" {
		t.Fatalf("description after update: got %q", got2.Description)
	}
	if got2.Mode.String() != "auth" {
		t.Fatalf("mode after update: got %q", got2.Mode.String())
	}
}
