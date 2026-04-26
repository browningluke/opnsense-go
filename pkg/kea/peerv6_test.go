package kea

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestPeerV6(t *testing.T) {
	api_client := api.NewClient(api.Options{
		Uri:           os.Getenv("OPNSENSE_URI"),
		APIKey:        os.Getenv("OPNSENSE_API_KEY"),
		APISecret:     os.Getenv("OPNSENSE_API_SECRET"),
		AllowInsecure: true,
		MaxBackoff:    30,
		MinBackoff:    1,
		MaxRetries:    4,
	})

	controller := Controller{Api: api_client}
	ctx := context.Background()

	peer := &PeerV6{
		Name: "test-kea-peer-v6",
		Url:  "http://[2001:db8::2]:647/",
		Role: api.SelectedMap("primary"),
	}

	key, err := controller.AddPeerV6(ctx, peer)
	if err != nil {
		t.Fatalf("Failed to add PeerV6: %v", err)
	}
	t.Logf("Added PeerV6 with key: %s", key)

	retrieved, err := controller.GetPeerV6(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get PeerV6: %v", err)
	}
	t.Logf("Retrieved PeerV6: %+v", retrieved)

	if retrieved.Name != peer.Name {
		t.Errorf("Name mismatch: got %s, want %s", retrieved.Name, peer.Name)
	}
	if retrieved.Url != peer.Url {
		t.Errorf("Url mismatch: got %s, want %s", retrieved.Url, peer.Url)
	}
	if retrieved.Role != peer.Role {
		t.Errorf("Role mismatch: got %s, want %s", retrieved.Role, peer.Role)
	}

	peer.Name = "test-kea-peer-v6-upd"
	peer.Url = "http://[2001:db8::3]:647/"
	err = controller.UpdatePeerV6(ctx, key, peer)
	if err != nil {
		t.Fatalf("Failed to update PeerV6: %v", err)
	}
	t.Logf("Updated PeerV6 with key: %s", key)

	retrieved, err = controller.GetPeerV6(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get updated PeerV6: %v", err)
	}
	if retrieved.Name != "test-kea-peer-v6-upd" {
		t.Errorf("Updated name mismatch: got %s, want %s", retrieved.Name, "test-kea-peer-v6-upd")
	}
	if retrieved.Url != "http://[2001:db8::3]:647/" {
		t.Errorf("Updated url mismatch: got %s, want %s", retrieved.Url, "http://[2001:db8::3]:647/")
	}

	err = controller.DeletePeerV6(ctx, key)
	if err != nil {
		t.Fatalf("Failed to delete PeerV6: %v", err)
	}
	t.Logf("Deleted PeerV6 with key: %s", key)
}
