package kea

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestReservationV6(t *testing.T) {
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

	// Create a SubnetV6 to associate the reservation with.
	subnet := &SubnetV6{
		Subnet:      "fd00:201::/64",
		Allocator:   api.SelectedMap("iterative"),
		PDAllocator: api.SelectedMap("iterative"),
		Pools:       "fd00:201::100-fd00:201::200",
		Interface:   api.SelectedMap("wan"),
		OptionData: OptionDataV6{
			DomainNameServers: api.SelectedMapList{},
			DomainSearch:      api.SelectedMapList{},
		},
		Description: "Test subnet for ReservationV6",
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

	reservation := &ReservationV6{
		Subnet:       api.SelectedMap(subnetKey),
		IpAddress:    "fd00:201::150",
		DUID:         "00:03:00:01:aa:bb:cc:dd:ee:ff",
		Hostname:     "test-host-v6",
		DomainSearch: api.SelectedMapList{},
		Description:  "Test Kea DHCPv6 Reservation",
	}

	key, err := controller.AddReservationV6(ctx, reservation)
	if err != nil {
		t.Fatalf("Failed to add ReservationV6: %v", err)
	}
	t.Logf("Added ReservationV6 with key: %s", key)

	retrieved, err := controller.GetReservationV6(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get ReservationV6: %v", err)
	}
	t.Logf("Retrieved ReservationV6: %+v", retrieved)

	if retrieved.Subnet != reservation.Subnet {
		t.Errorf("Subnet mismatch: got %s, want %s", retrieved.Subnet, reservation.Subnet)
	}
	if retrieved.IpAddress != reservation.IpAddress {
		t.Errorf("IpAddress mismatch: got %s, want %s", retrieved.IpAddress, reservation.IpAddress)
	}
	if retrieved.DUID != reservation.DUID {
		t.Errorf("DUID mismatch: got %s, want %s", retrieved.DUID, reservation.DUID)
	}
	if retrieved.Hostname != reservation.Hostname {
		t.Errorf("Hostname mismatch: got %s, want %s", retrieved.Hostname, reservation.Hostname)
	}
	if retrieved.Description != reservation.Description {
		t.Errorf("Description mismatch: got %s, want %s", retrieved.Description, reservation.Description)
	}

	reservation.Hostname = "test-host-v6-upd"
	reservation.Description = "Test Kea DHCPv6 Reservation Updated"
	err = controller.UpdateReservationV6(ctx, key, reservation)
	if err != nil {
		t.Fatalf("Failed to update ReservationV6: %v", err)
	}
	t.Logf("Updated ReservationV6 with key: %s", key)

	retrieved, err = controller.GetReservationV6(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get updated ReservationV6: %v", err)
	}
	if retrieved.Hostname != "test-host-v6-upd" {
		t.Errorf("Updated hostname mismatch: got %s, want %s", retrieved.Hostname, "test-host-v6-upd")
	}
	if retrieved.Description != "Test Kea DHCPv6 Reservation Updated" {
		t.Errorf("Updated description mismatch: got %s, want %s", retrieved.Description, "Test Kea DHCPv6 Reservation Updated")
	}

	err = controller.DeleteReservationV6(ctx, key)
	if err != nil {
		t.Fatalf("Failed to delete ReservationV6: %v", err)
	}
	t.Logf("Deleted ReservationV6 with key: %s", key)
}
