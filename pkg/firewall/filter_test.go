package firewall

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestFilter(t *testing.T) {
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

	ctx := context.Background()

	controller := Controller{
		Api: api_client,
	}

	filter := &Filter{
		Enabled:                  "0",
		Sequence:                 "1",
		Categories:               api.SelectedMapList{},
		NoXMLRPCSync:             "0",
		Description:              "Test filter rule",
		InvertInterface:          "0",
		Interface:                api.SelectedMapList{"wan"},
		Quick:                    "1",
		Action:                   api.SelectedMap("pass"),
		AllowOptions:             "0",
		Direction:                api.SelectedMap("in"),
		IPProtocol:               api.SelectedMap("inet"),
		Protocol:                 api.SelectedMap("TCP"),
		ICMPType:                 api.SelectedMapList{},
		SourceInvert:             "0",
		SourceNet:                "any",
		SourcePort:               "",
		DestinationInvert:        "0",
		DestinationNet:           "192.168.1.0/24",
		DestinationPort:          "80",
		Log:                      "0",
		TCPFlags:                 api.SelectedMapList{"syn"},
		TCPFlagsOutOf:            api.SelectedMapList{"syn", "ack"},
		Schedule:                 api.SelectedMap(""),
		StateType:                api.SelectedMap("keep"),
		StatePolicy:              api.SelectedMap(""),
		StateTimeout:             "3600",
		AdaptiveTimeoutsStart:    "1000",
		AdaptiveTimeoutsEnd:      "15000",
		MaxStates:                "10000",
		MaxSourceNodes:           "100",
		MaxSourceStates:          "500",
		MaxSourceConnections:     "50",
		MaxNewConnectionsCount:   "10",
		MaxNewConnectionsSeconds: "60",
		OverloadTable:            api.SelectedMap(""),
		NoPfsync:                 "0",
		TrafficShaper:            api.SelectedMap(""),
		TrafficShaperReverse:     api.SelectedMap(""),
		Gateway:                  api.SelectedMap(""),
		DisableReplyTo:           "0",
		ReplyTo:                  api.SelectedMap(""),
		MatchPriority:            api.SelectedMap(""),
		SetPriority:              api.SelectedMap(""),
		SetPriorityLowDelay:      api.SelectedMap(""),
		MatchTOS:                 api.SelectedMap(""),
		SetLocalTag:              "",
		MatchLocalTag:            "",
	}

	key, err := controller.AddFilter(ctx, filter)
	if err != nil {
		t.Fatalf("Failed to add filter rule: %v", err)
	}
	t.Logf("Added filter rule with key: %s", key)

	retrievedFilter, err := controller.GetFilter(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get filter rule: %v", err)
	}
	t.Logf("Retrieved filter rule: %+v", retrievedFilter)

	// Test original fields
	if retrievedFilter.SourceNet != filter.SourceNet {
		t.Fatalf("Retrieved filter source net does not match: got %s, want %s", retrievedFilter.SourceNet, filter.SourceNet)
	}
	if retrievedFilter.DestinationNet != filter.DestinationNet {
		t.Fatalf("Retrieved filter destination net does not match: got %s, want %s", retrievedFilter.DestinationNet, filter.DestinationNet)
	}
	if retrievedFilter.Description != filter.Description {
		t.Fatalf("Retrieved filter description does not match: got %s, want %s", retrievedFilter.Description, filter.Description)
	}

	// Test new fields
	if retrievedFilter.Enabled != filter.Enabled {
		t.Fatalf("Retrieved filter enabled does not match: got %s, want %s", retrievedFilter.Enabled, filter.Enabled)
	}
	if retrievedFilter.NoXMLRPCSync != filter.NoXMLRPCSync {
		t.Fatalf("Retrieved filter NoXMLRPCSync does not match: got %s, want %s", retrievedFilter.NoXMLRPCSync, filter.NoXMLRPCSync)
	}
	if retrievedFilter.InvertInterface != filter.InvertInterface {
		t.Fatalf("Retrieved filter InvertInterface does not match: got %s, want %s", retrievedFilter.InvertInterface, filter.InvertInterface)
	}
	if retrievedFilter.AllowOptions != filter.AllowOptions {
		t.Fatalf("Retrieved filter AllowOptions does not match: got %s, want %s", retrievedFilter.AllowOptions, filter.AllowOptions)
	}
	if retrievedFilter.StateType != filter.StateType {
		t.Fatalf("Retrieved filter StateType does not match: got %s, want %s", retrievedFilter.StateType, filter.StateType)
	}
	if retrievedFilter.StateTimeout != filter.StateTimeout {
		t.Fatalf("Retrieved filter StateTimeout does not match: got %s, want %s", retrievedFilter.StateTimeout, filter.StateTimeout)
	}
	if retrievedFilter.AdaptiveTimeoutsStart != filter.AdaptiveTimeoutsStart {
		t.Fatalf("Retrieved filter AdaptiveTimeoutsStart does not match: got %s, want %s", retrievedFilter.AdaptiveTimeoutsStart, filter.AdaptiveTimeoutsStart)
	}
	if retrievedFilter.AdaptiveTimeoutsEnd != filter.AdaptiveTimeoutsEnd {
		t.Fatalf("Retrieved filter AdaptiveTimeoutsEnd does not match: got %s, want %s", retrievedFilter.AdaptiveTimeoutsEnd, filter.AdaptiveTimeoutsEnd)
	}
	if retrievedFilter.MaxStates != filter.MaxStates {
		t.Fatalf("Retrieved filter MaxStates does not match: got %s, want %s", retrievedFilter.MaxStates, filter.MaxStates)
	}
	if retrievedFilter.MaxSourceNodes != filter.MaxSourceNodes {
		t.Fatalf("Retrieved filter MaxSourceNodes does not match: got %s, want %s", retrievedFilter.MaxSourceNodes, filter.MaxSourceNodes)
	}
	if retrievedFilter.MaxSourceStates != filter.MaxSourceStates {
		t.Fatalf("Retrieved filter MaxSourceStates does not match: got %s, want %s", retrievedFilter.MaxSourceStates, filter.MaxSourceStates)
	}
	if retrievedFilter.MaxSourceConnections != filter.MaxSourceConnections {
		t.Fatalf("Retrieved filter MaxSourceConnections does not match: got %s, want %s", retrievedFilter.MaxSourceConnections, filter.MaxSourceConnections)
	}
	if retrievedFilter.MaxNewConnectionsCount != filter.MaxNewConnectionsCount {
		t.Fatalf("Retrieved filter MaxNewConnectionsCount does not match: got %s, want %s", retrievedFilter.MaxNewConnectionsCount, filter.MaxNewConnectionsCount)
	}
	if retrievedFilter.MaxNewConnectionsSeconds != filter.MaxNewConnectionsSeconds {
		t.Fatalf("Retrieved filter MaxNewConnectionsSeconds does not match: got %s, want %s", retrievedFilter.MaxNewConnectionsSeconds, filter.MaxNewConnectionsSeconds)
	}
	if retrievedFilter.NoPfsync != filter.NoPfsync {
		t.Fatalf("Retrieved filter NoPfsync does not match: got %s, want %s", retrievedFilter.NoPfsync, filter.NoPfsync)
	}
	if retrievedFilter.DisableReplyTo != filter.DisableReplyTo {
		t.Fatalf("Retrieved filter DisableReplyTo does not match: got %s, want %s", retrievedFilter.DisableReplyTo, filter.DisableReplyTo)
	}
	if len(retrievedFilter.TCPFlags) != len(filter.TCPFlags) {
		t.Fatalf("Retrieved filter TCPFlags length does not match: got %d, want %d", len(retrievedFilter.TCPFlags), len(filter.TCPFlags))
	}
	if len(retrievedFilter.TCPFlagsOutOf) != len(filter.TCPFlagsOutOf) {
		t.Fatalf("Retrieved filter TCPFlagsOutOf length does not match: got %d, want %d", len(retrievedFilter.TCPFlagsOutOf), len(filter.TCPFlagsOutOf))
	}

	filter.DestinationNet = "192.168.2.0/24"
	filter.DestinationPort = "443"
	filter.Description = "Test filter rule updated"
	filter.Enabled = "1"
	filter.Log = "1"
	filter.StateTimeout = "7200"
	filter.AdaptiveTimeoutsEnd = "25000"
	filter.MaxStates = "20000"
	filter.MaxSourceNodes = "200"
	filter.MaxNewConnectionsCount = "20"
	filter.TCPFlags = api.SelectedMapList{"syn", "ack"}
	err = controller.UpdateFilter(ctx, key, filter)
	if err != nil {
		t.Fatalf("Failed to update filter rule: %v", err)
	}
	t.Logf("Updated filter rule with key: %s", key)

	retrievedFilter, err = controller.GetFilter(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get updated filter rule: %v", err)
	}
	if retrievedFilter.DestinationNet != "192.168.2.0/24" {
		t.Fatalf("Retrieved filter destination net does not match updated net: got %s, want %s", retrievedFilter.DestinationNet, "192.168.2.0/24")
	}
	if retrievedFilter.DestinationPort != "443" {
		t.Fatalf("Retrieved filter destination port does not match updated port: got %s, want %s", retrievedFilter.DestinationPort, "443")
	}
	if retrievedFilter.Description != "Test filter rule updated" {
		t.Fatalf("Retrieved filter description does not match updated description: got %s, want %s", retrievedFilter.Description, "Test filter rule updated")
	}
	if retrievedFilter.Enabled != "1" {
		t.Fatalf("Retrieved filter enabled does not match updated value: got %s, want %s", retrievedFilter.Enabled, "1")
	}
	if retrievedFilter.Log != "1" {
		t.Fatalf("Retrieved filter log does not match updated value: got %s, want %s", retrievedFilter.Log, "1")
	}
	if retrievedFilter.StateTimeout != "7200" {
		t.Fatalf("Retrieved filter StateTimeout does not match updated value: got %s, want %s", retrievedFilter.StateTimeout, "7200")
	}
	if retrievedFilter.AdaptiveTimeoutsEnd != "25000" {
		t.Fatalf("Retrieved filter AdaptiveTimeoutsEnd does not match updated value: got %s, want %s", retrievedFilter.AdaptiveTimeoutsEnd, "25000")
	}
	if retrievedFilter.MaxStates != "20000" {
		t.Fatalf("Retrieved filter MaxStates does not match updated value: got %s, want %s", retrievedFilter.MaxStates, "20000")
	}
	if retrievedFilter.MaxSourceNodes != "200" {
		t.Fatalf("Retrieved filter MaxSourceNodes does not match updated value: got %s, want %s", retrievedFilter.MaxSourceNodes, "200")
	}
	if retrievedFilter.MaxNewConnectionsCount != "20" {
		t.Fatalf("Retrieved filter MaxNewConnectionsCount does not match updated value: got %s, want %s", retrievedFilter.MaxNewConnectionsCount, "20")
	}
	if len(retrievedFilter.TCPFlags) != 2 {
		t.Fatalf("Retrieved filter TCPFlags length does not match updated value: got %d, want %d", len(retrievedFilter.TCPFlags), 2)
	}

	err = controller.DeleteFilter(ctx, key)
	if err != nil {
		t.Fatalf("Failed to delete filter rule: %v", err)
	}
	t.Logf("Deleted filter rule with key: %s", key)
}
