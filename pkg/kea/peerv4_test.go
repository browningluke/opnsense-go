package kea

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestPeerV4(t *testing.T) {
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

	peer := &PeerV4{
		Name: "test-kea-peer-v4",
		Url:  "http://192.168.1.2:647/",
		Role: api.SelectedMap("primary"),
	}

	key, err := controller.AddPeerV4(ctx, peer)
	if err != nil {
		t.Fatalf("Failed to add PeerV4: %v", err)
	}
	t.Logf("Added PeerV4 with key: %s", key)

	retrieved, err := controller.GetPeerV4(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get PeerV4: %v", err)
	}
	t.Logf("Retrieved PeerV4: %+v", retrieved)

	if retrieved.Name != peer.Name {
		t.Errorf("Name mismatch: got %s, want %s", retrieved.Name, peer.Name)
	}
	if retrieved.Url != peer.Url {
		t.Errorf("Url mismatch: got %s, want %s", retrieved.Url, peer.Url)
	}
	if retrieved.Role != peer.Role {
		t.Errorf("Role mismatch: got %s, want %s", retrieved.Role, peer.Role)
	}

	peer.Name = "test-kea-peer-v4-upd"
	peer.Url = "http://192.168.1.3:647/"
	err = controller.UpdatePeerV4(ctx, key, peer)
	if err != nil {
		t.Fatalf("Failed to update PeerV4: %v", err)
	}
	t.Logf("Updated PeerV4 with key: %s", key)

	retrieved, err = controller.GetPeerV4(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get updated PeerV4: %v", err)
	}
	if retrieved.Name != "test-kea-peer-v4-upd" {
		t.Errorf("Updated name mismatch: got %s, want %s", retrieved.Name, "test-kea-peer-v4-upd")
	}
	if retrieved.Url != "http://192.168.1.3:647/" {
		t.Errorf("Updated url mismatch: got %s, want %s", retrieved.Url, "http://192.168.1.3:647/")
	}

	err = controller.DeletePeerV4(ctx, key)
	if err != nil {
		t.Fatalf("Failed to delete PeerV4: %v", err)
	}
	t.Logf("Deleted PeerV4 with key: %s", key)
}
