package kea

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestReservationV4(t *testing.T) {
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

	// Create a SubnetV4 to associate the reservation with.
	subnet := &SubnetV4{
		Subnet:                "192.168.201.0/24",
		Pools:                 "192.168.201.100 - 192.168.201.200",
		MatchClientId:         "1",
		OptionDataAutoCollect: "1",
		OptionData: OptionDataV4{
			DomainNameServers: api.SelectedMapList{},
			Routers:           api.SelectedMapList{"192.168.201.1"},
			NtpServers:        api.SelectedMapList{},
			TimeServers:       api.SelectedMapList{},
			DomainSearch:      api.SelectedMapList{},
		},
		Description: "Test subnet for ReservationV4",
	}
	subnetKey, err := controller.AddSubnetV4(ctx, subnet)
	if err != nil {
		t.Fatalf("Failed to add SubnetV4 fixture: %v", err)
	}
	t.Logf("Created SubnetV4 fixture with key: %s", subnetKey)
	defer func() {
		if delErr := controller.DeleteSubnetV4(ctx, subnetKey); delErr != nil {
			t.Errorf("Failed to delete SubnetV4 fixture: %v", delErr)
		} else {
			t.Logf("Deleted SubnetV4 fixture with key: %s", subnetKey)
		}
	}()

	reservation := &ReservationV4{
		Subnet:      api.SelectedMap(subnetKey),
		IpAddress:   "192.168.201.150",
		HwAddress:   "aa:bb:cc:dd:ee:ff",
		Hostname:    "test-host",
		Description: "Test Kea DHCPv4 Reservation",
	}

	key, err := controller.AddReservationV4(ctx, reservation)
	if err != nil {
		t.Fatalf("Failed to add ReservationV4: %v", err)
	}
	t.Logf("Added ReservationV4 with key: %s", key)

	retrieved, err := controller.GetReservationV4(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get ReservationV4: %v", err)
	}
	t.Logf("Retrieved ReservationV4: %+v", retrieved)

	if retrieved.Subnet != reservation.Subnet {
		t.Errorf("Subnet mismatch: got %s, want %s", retrieved.Subnet, reservation.Subnet)
	}
	if retrieved.IpAddress != reservation.IpAddress {
		t.Errorf("IpAddress mismatch: got %s, want %s", retrieved.IpAddress, reservation.IpAddress)
	}
	if retrieved.HwAddress != reservation.HwAddress {
		t.Errorf("HwAddress mismatch: got %s, want %s", retrieved.HwAddress, reservation.HwAddress)
	}
	if retrieved.Hostname != reservation.Hostname {
		t.Errorf("Hostname mismatch: got %s, want %s", retrieved.Hostname, reservation.Hostname)
	}
	if retrieved.Description != reservation.Description {
		t.Errorf("Description mismatch: got %s, want %s", retrieved.Description, reservation.Description)
	}

	reservation.Hostname = "test-host-upd"
	reservation.Description = "Test Kea DHCPv4 Reservation Updated"
	err = controller.UpdateReservationV4(ctx, key, reservation)
	if err != nil {
		t.Fatalf("Failed to update ReservationV4: %v", err)
	}
	t.Logf("Updated ReservationV4 with key: %s", key)

	retrieved, err = controller.GetReservationV4(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get updated ReservationV4: %v", err)
	}
	if retrieved.Hostname != "test-host-upd" {
		t.Errorf("Updated hostname mismatch: got %s, want %s", retrieved.Hostname, "test-host-upd")
	}
	if retrieved.Description != "Test Kea DHCPv4 Reservation Updated" {
		t.Errorf("Updated description mismatch: got %s, want %s", retrieved.Description, "Test Kea DHCPv4 Reservation Updated")
	}

	err = controller.DeleteReservationV4(ctx, key)
	if err != nil {
		t.Fatalf("Failed to delete ReservationV4: %v", err)
	}
	t.Logf("Deleted ReservationV4 with key: %s", key)
}
