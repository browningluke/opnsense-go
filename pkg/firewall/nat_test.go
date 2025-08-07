package firewall

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestNAT(t *testing.T) {
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

	nat := &NAT{
		Enabled:           "1",
		DisableNAT:        "0",
		Sequence:          "1",
		Interface:         api.SelectedMap("wan"),
		IPProtocol:        api.SelectedMap("inet"),
		Protocol:          api.SelectedMap("TCP"),
		SourceNet:         "any",
		SourcePort:        "",
		SourceInvert:      "0",
		DestinationNet:    "any",
		DestinationPort:   "80",
		DestinationInvert: "0",
		Target:            "192.168.1.100",
		TargetPort:        "8080",
		Log:               "0",
		Description:       "Test NAT rule",
	}

	key, err := controller.AddNAT(ctx, nat)
	if err != nil {
		t.Fatalf("Failed to add NAT rule: %v", err)
	}
	t.Logf("Added NAT rule with key: %s", key)

	retrievedNAT, err := controller.GetNAT(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get NAT rule: %v", err)
	}
	t.Logf("Retrieved NAT rule: %+v", retrievedNAT)
	if retrievedNAT.SourceNet != nat.SourceNet {
		t.Fatalf("Retrieved NAT source net does not match: got %s, want %s", retrievedNAT.SourceNet, nat.SourceNet)
	}
	if retrievedNAT.Target != nat.Target {
		t.Fatalf("Retrieved NAT target does not match: got %s, want %s", retrievedNAT.Target, nat.Target)
	}
	if retrievedNAT.Description != nat.Description {
		t.Fatalf("Retrieved NAT description does not match: got %s, want %s", retrievedNAT.Description, nat.Description)
	}

	nat.Target = "192.168.1.200"
	nat.TargetPort = "9090"
	nat.Description = "Test NAT rule updated"
	err = controller.UpdateNAT(ctx, key, nat)
	if err != nil {
		t.Fatalf("Failed to update NAT rule: %v", err)
	}
	t.Logf("Updated NAT rule with key: %s", key)

	retrievedNAT, err = controller.GetNAT(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get updated NAT rule: %v", err)
	}
	if retrievedNAT.Target != "192.168.1.200" {
		t.Fatalf("Retrieved NAT target does not match updated target: got %s, want %s", retrievedNAT.Target, "192.168.1.200")
	}
	if retrievedNAT.TargetPort != "9090" {
		t.Fatalf("Retrieved NAT target port does not match updated target port: got %s, want %s", retrievedNAT.TargetPort, "9090")
	}
	if retrievedNAT.Description != "Test NAT rule updated" {
		t.Fatalf("Retrieved NAT description does not match updated description: got %s, want %s", retrievedNAT.Description, "Test NAT rule updated")
	}

	err = controller.DeleteNAT(ctx, key)
	if err != nil {
		t.Fatalf("Failed to delete NAT rule: %v", err)
	}
	t.Logf("Deleted NAT rule with key: %s", key)
}
