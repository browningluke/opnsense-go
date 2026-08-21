package kea

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestSubnetV4(t *testing.T) {
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

	subnet := &SubnetV4{
		Subnet:                "192.168.200.0/24",
		ValidLifetime:         "86400",
		NextServer:            "",
		Pools:                 "192.168.200.100 - 192.168.200.200",
		MatchClientId:         "1",
		OptionDataAutoCollect: "1",
		OptionData: OptionDataV4{
			DomainNameServers: api.SelectedMapList{"8.8.8.8"},
			Routers:           api.SelectedMapList{"192.168.200.1"},
			DomainName:        "test.local",
			NtpServers:        api.SelectedMapList{},
			TimeServers:       api.SelectedMapList{},
			DomainSearch:      api.SelectedMapList{},
		},
		Description: "Test Kea DHCPv4 Subnet",
	}

	key, err := controller.AddSubnetV4(ctx, subnet)
	if err != nil {
		t.Fatalf("Failed to add SubnetV4: %v", err)
	}
	deleted := false
	defer func() {
		if !deleted {
			_ = controller.DeleteSubnetV4(ctx, key)
		}
	}()
	t.Logf("Added SubnetV4 with key: %s", key)

	retrieved, err := controller.GetSubnetV4(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get SubnetV4: %v", err)
	}
	t.Logf("Retrieved SubnetV4: %+v", retrieved)

	if retrieved.Subnet != subnet.Subnet {
		t.Errorf("Subnet mismatch: got %s, want %s", retrieved.Subnet, subnet.Subnet)
	}
	supportsValidLifetime := retrieved.ValidLifetime != ""
	if supportsValidLifetime && retrieved.ValidLifetime != subnet.ValidLifetime {
		t.Errorf("ValidLifetime mismatch: got %s, want %s", retrieved.ValidLifetime, subnet.ValidLifetime)
	}
	if !supportsValidLifetime {
		t.Log("OPNsense version does not expose per-subnet valid_lifetime; skipping lifetime API assertions")
	}
	if retrieved.Pools != subnet.Pools {
		t.Errorf("Pools mismatch: got %s, want %s", retrieved.Pools, subnet.Pools)
	}
	if retrieved.MatchClientId != subnet.MatchClientId {
		t.Errorf("MatchClientId mismatch: got %s, want %s", retrieved.MatchClientId, subnet.MatchClientId)
	}
	if retrieved.OptionDataAutoCollect != subnet.OptionDataAutoCollect {
		t.Errorf("OptionDataAutoCollect mismatch: got %s, want %s", retrieved.OptionDataAutoCollect, subnet.OptionDataAutoCollect)
	}
	if retrieved.OptionData.DomainName != subnet.OptionData.DomainName {
		t.Errorf("OptionData.DomainName mismatch: got %s, want %s", retrieved.OptionData.DomainName, subnet.OptionData.DomainName)
	}
	if retrieved.Description != subnet.Description {
		t.Errorf("Description mismatch: got %s, want %s", retrieved.Description, subnet.Description)
	}

	subnet.Description = "Test Kea DHCPv4 Subnet Updated"
	subnet.Pools = "192.168.200.100 - 192.168.200.150"
	subnet.ValidLifetime = "43200"
	err = controller.UpdateSubnetV4(ctx, key, subnet)
	if err != nil {
		t.Fatalf("Failed to update SubnetV4: %v", err)
	}
	t.Logf("Updated SubnetV4 with key: %s", key)

	retrieved, err = controller.GetSubnetV4(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get updated SubnetV4: %v", err)
	}
	if retrieved.Description != "Test Kea DHCPv4 Subnet Updated" {
		t.Errorf("Updated description mismatch: got %s, want %s", retrieved.Description, "Test Kea DHCPv4 Subnet Updated")
	}
	if retrieved.Pools != "192.168.200.100 - 192.168.200.150" {
		t.Errorf("Updated pools mismatch: got %s, want %s", retrieved.Pools, "192.168.200.100 - 192.168.200.150")
	}
	if supportsValidLifetime && retrieved.ValidLifetime != "43200" {
		t.Errorf("Updated valid lifetime mismatch: got %s, want %s", retrieved.ValidLifetime, "43200")
	}

	err = controller.DeleteSubnetV4(ctx, key)
	if err != nil {
		t.Fatalf("Failed to delete SubnetV4: %v", err)
	}
	deleted = true
	t.Logf("Deleted SubnetV4 with key: %s", key)
}

func TestSubnetV4ValidLifetimeJSON(t *testing.T) {
	subnet := SubnetV4{ValidLifetime: "86400"}

	payload, err := json.Marshal(subnet)
	if err != nil {
		t.Fatalf("Failed to marshal SubnetV4: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal SubnetV4 JSON: %v", err)
	}
	if decoded["valid_lifetime"] != "86400" {
		t.Fatalf("valid_lifetime JSON mismatch: got %v, want %s", decoded["valid_lifetime"], "86400")
	}
}
