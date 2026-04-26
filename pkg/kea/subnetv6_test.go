package kea

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestSubnetV6(t *testing.T) {
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

	subnet := &SubnetV6{
		Subnet:      "fd00:200::/64",
		Allocator:   api.SelectedMap("iterative"),
		PDAllocator: api.SelectedMap("iterative"),
		Pools:       "fd00:200::100-fd00:200::200",
		Interface:   api.SelectedMap("wan"),
		OptionData: OptionDataV6{
			DomainNameServers: api.SelectedMapList{"2001:4860:4860::8888"},
			DomainSearch:      api.SelectedMapList{},
		},
		Description: "Test Kea DHCPv6 Subnet",
	}

	key, err := controller.AddSubnetV6(ctx, subnet)
	if err != nil {
		t.Fatalf("Failed to add SubnetV6: %v", err)
	}
	t.Logf("Added SubnetV6 with key: %s", key)

	retrieved, err := controller.GetSubnetV6(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get SubnetV6: %v", err)
	}
	t.Logf("Retrieved SubnetV6: %+v", retrieved)

	if retrieved.Subnet != subnet.Subnet {
		t.Errorf("Subnet mismatch: got %s, want %s", retrieved.Subnet, subnet.Subnet)
	}
	if retrieved.Allocator != subnet.Allocator {
		t.Errorf("Allocator mismatch: got %s, want %s", retrieved.Allocator, subnet.Allocator)
	}
	if retrieved.PDAllocator != subnet.PDAllocator {
		t.Errorf("PDAllocator mismatch: got %s, want %s", retrieved.PDAllocator, subnet.PDAllocator)
	}
	if retrieved.Pools != subnet.Pools {
		t.Errorf("Pools mismatch: got %s, want %s", retrieved.Pools, subnet.Pools)
	}
	if retrieved.Description != subnet.Description {
		t.Errorf("Description mismatch: got %s, want %s", retrieved.Description, subnet.Description)
	}

	subnet.Description = "Test Kea DHCPv6 Subnet Updated"
	subnet.Pools = "fd00:200::100-fd00:200::150"
	err = controller.UpdateSubnetV6(ctx, key, subnet)
	if err != nil {
		t.Fatalf("Failed to update SubnetV6: %v", err)
	}
	t.Logf("Updated SubnetV6 with key: %s", key)

	retrieved, err = controller.GetSubnetV6(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get updated SubnetV6: %v", err)
	}
	if retrieved.Description != "Test Kea DHCPv6 Subnet Updated" {
		t.Errorf("Updated description mismatch: got %s, want %s", retrieved.Description, "Test Kea DHCPv6 Subnet Updated")
	}
	if retrieved.Pools != "fd00:200::100-fd00:200::150" {
		t.Errorf("Updated pools mismatch: got %s, want %s", retrieved.Pools, "fd00:200::100-fd00:200::150")
	}

	err = controller.DeleteSubnetV6(ctx, key)
	if err != nil {
		t.Fatalf("Failed to delete SubnetV6: %v", err)
	}
	t.Logf("Deleted SubnetV6 with key: %s", key)
}
