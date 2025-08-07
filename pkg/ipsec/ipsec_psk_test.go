package ipsec

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestPSK(t *testing.T) {
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

	ipsec_psk := &IPsecPSK{
		IdentityLocal:  "my-test-uuid",
		IdentityRemote: "my-test-remote-uuid",
		PreSharedKey:   "my-test-psk",
		Type:           "PSK",
		Description:    "Test PSK",
	}
	key, err := controller.AddIPsecPSK(ctx, ipsec_psk)
	if err != nil {
		t.Fatalf("Failed to add IPsec PSK: %v", err)
	}
	t.Logf("Added IPsec PSK with key: %s", key)

	rekey, err := controller.GetIPsecPSK(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get IPsec PSK: %v", err)
	}
	t.Logf("Retrieved IPsec PSK: %+v", rekey)
	if rekey.IdentityLocal != ipsec_psk.IdentityLocal || rekey.IdentityRemote != ipsec_psk.IdentityRemote {
		t.Fatalf("Retrieved PSK does not match added PSK: got %+v, want %+v", rekey, ipsec_psk)
	}
	if rekey.PreSharedKey != ipsec_psk.PreSharedKey {
		t.Fatalf("Retrieved PSK key does not match added PSK key: got %s, want %s", rekey.PreSharedKey, ipsec_psk.PreSharedKey)
	}
	if rekey.Type != ipsec_psk.Type {
		t.Fatalf("Retrieved PSK type does not match added PSK type: got %s, want %s", rekey.Type, ipsec_psk.Type)
	}
	if rekey.Description != ipsec_psk.Description {
		t.Fatalf("Retrieved PSK description does not match added PSK description: got %s, want %s", rekey.Description, ipsec_psk.Description)
	}
	t.Logf("Successfully verified retrieved PSK matches added PSK")

	ipsec_psk.IdentityLocal = "my-test-uuid-updated"
	ipsec_psk.IdentityRemote = "my-test-remote-uuid-updated"
	ipsec_psk.PreSharedKey = "my-test-psk-updated"
	ipsec_psk.Type = "PSK"
	ipsec_psk.Description = "Test PSK Updated"
	t.Logf("Updating IPsec PSK with key: %s", key)

	err = controller.UpdateIPsecPSK(ctx, key, ipsec_psk)
	if err != nil {
		t.Fatalf("Failed to update IPsec PSK: %v", err)
	}
	t.Logf("Updated IPsec PSK with key: %s", key)
	rekey, err = controller.GetIPsecPSK(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get updated IPsec PSK: %v", err)
	}
	t.Logf("Retrieved updated IPsec PSK: %+v", rekey)
	if rekey.IdentityLocal != "my-test-uuid-updated" || rekey.IdentityRemote != "my-test-remote-uuid-updated" {
		t.Fatalf("Retrieved PSK does not match updated PSK: got %+v, want %+v", rekey, ipsec_psk)
	}
	if rekey.PreSharedKey != "my-test-psk-updated" {
		t.Fatalf("Retrieved PSK key does not match updated PSK key: got %s, want %s", rekey.PreSharedKey, "my-test-psk-updated")
	}
	if rekey.Type != "PSK" {
		t.Fatalf("Retrieved PSK type does not match updated PSK type: got %s, want %s", rekey.Type, "PSK")
	}
	if rekey.Description != "Test PSK Updated" {
		t.Fatalf("Retrieved PSK description does not match updated PSK description: got %s, want %s", rekey.Description, "Test PSK Updated")
	}
	t.Logf("Successfully verified retrieved PSK matches updated PSK")

	err = controller.DeleteIPsecPSK(ctx, key)
	if err != nil {
		t.Fatalf("Failed to delete IPsec PSK: %v", err)
	}
	t.Logf("Deleted IPsec PSK with key: %s", key)
}
