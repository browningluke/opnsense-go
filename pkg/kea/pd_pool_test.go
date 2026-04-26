package kea

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestPDPool(t *testing.T) {
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

	// Create a SubnetV6 to associate the PD pool with.
	subnet := &SubnetV6{
		Subnet:      "fd00:202::/64",
		Allocator:   api.SelectedMap("iterative"),
		PDAllocator: api.SelectedMap("iterative"),
		Pools:       "",
		Interface:   api.SelectedMap("wan"),
		OptionData: OptionDataV6{
			DomainNameServers: api.SelectedMapList{},
			DomainSearch:      api.SelectedMapList{},
		},
		Description: "Test subnet for PDPool",
	}
	subnetKey, err := controller.AddSubnetV6(ctx, subnet)
	if err != nil {
		t.Fatalf("Failed to add SubnetV6 fixture: %v", err)
	}
	t.Logf("Created SubnetV6 fixture with key: %s", subnetKey)
	defer func() {
		if delErr := controller.DeleteSubnetV6(ctx, subnetKey); delErr != nil {
			t.Errorf("Failed to delete SubnetV6 fixture: %v", delErr)
		} else {
			t.Logf("Deleted SubnetV6 fixture with key: %s", subnetKey)
		}
	}()

	pdpool := &PDPool{
		Subnet:       api.SelectedMap(subnetKey),
		Prefix:       "fd00:203::",
		PrefixLen:    "64",
		DelegatedLen: "80",
		Description:  "Test Kea PD Pool",
	}

	key, err := controller.AddPDPool(ctx, pdpool)
	if err != nil {
		t.Fatalf("Failed to add PDPool: %v", err)
	}
	t.Logf("Added PDPool with key: %s", key)

	retrieved, err := controller.GetPDPool(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get PDPool: %v", err)
	}
	t.Logf("Retrieved PDPool: %+v", retrieved)

	if retrieved.Subnet != pdpool.Subnet {
		t.Errorf("Subnet mismatch: got %s, want %s", retrieved.Subnet, pdpool.Subnet)
	}
	if retrieved.Prefix != pdpool.Prefix {
		t.Errorf("Prefix mismatch: got %s, want %s", retrieved.Prefix, pdpool.Prefix)
	}
	if retrieved.PrefixLen != pdpool.PrefixLen {
		t.Errorf("PrefixLen mismatch: got %s, want %s", retrieved.PrefixLen, pdpool.PrefixLen)
	}
	if retrieved.DelegatedLen != pdpool.DelegatedLen {
		t.Errorf("DelegatedLen mismatch: got %s, want %s", retrieved.DelegatedLen, pdpool.DelegatedLen)
	}
	if retrieved.Description != pdpool.Description {
		t.Errorf("Description mismatch: got %s, want %s", retrieved.Description, pdpool.Description)
	}

	pdpool.Description = "Test Kea PD Pool Updated"
	pdpool.DelegatedLen = "96"
	err = controller.UpdatePDPool(ctx, key, pdpool)
	if err != nil {
		t.Fatalf("Failed to update PDPool: %v", err)
	}
	t.Logf("Updated PDPool with key: %s", key)

	retrieved, err = controller.GetPDPool(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get updated PDPool: %v", err)
	}
	if retrieved.Description != "Test Kea PD Pool Updated" {
		t.Errorf("Updated description mismatch: got %s, want %s", retrieved.Description, "Test Kea PD Pool Updated")
	}
	if retrieved.DelegatedLen != "96" {
		t.Errorf("Updated delegated_len mismatch: got %s, want %s", retrieved.DelegatedLen, "96")
	}

	err = controller.DeletePDPool(ctx, key)
	if err != nil {
		t.Fatalf("Failed to delete PDPool: %v", err)
	}
	t.Logf("Deleted PDPool with key: %s", key)
}
