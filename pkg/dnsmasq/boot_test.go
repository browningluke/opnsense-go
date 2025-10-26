package dnsmasq

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestBoot(t *testing.T) {
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

	boot := &Boot{
		Interface: api.SelectedMap("lan"),
		// Tag:           api.SelectedMapList([]string{"b4b79319-2e09-47ae-85c2-687eb4b6e7ee", "d594fa8a-1a76-44e7-afab-adb6c5bdb69e"}),
		Filename:      "test",
		Servername:    "test-servername",
		Serveraddress: "test-serveraddress",
		Description:   "test-description",
	}

	respAdd, err := controller.AddBoot(ctx, boot)
	if err != nil {
		t.Fatalf("Failed to add boot: %v", err)
	}
	t.Logf("AddBoot: %+v", respAdd)

	respGet, err := controller.GetBoot(ctx, respAdd)
	if err != nil {
		t.Fatalf("Failed to get boot: %v", err)
	}
	t.Logf("GetBoot: %+v", respGet)

	boot.Interface = api.SelectedMap("wan")
	// boot.Tag = api.SelectedMapList([]string{"b4b79319-2e09-47ae-85c2-687eb4b6e7ee"})
	boot.Filename = "test-updated"
	boot.Servername = "test-servername-updated"
	boot.Serveraddress = "test-serveraddress-updated"
	boot.Description = "test-description-updated"
	err = controller.UpdateBoot(ctx, respAdd, boot)
	if err != nil {
		t.Fatalf("Failed to update boot: %v", err)
	}
	t.Logf("UpdateBoot: %+v", boot)

	respGet, err = controller.GetBoot(ctx, respAdd)
	if err != nil {
		t.Fatalf("Failed to get boot: %v", err)
	}
	t.Logf("GetBoot: %+v", respGet)

	err = controller.DeleteBoot(ctx, respAdd)
	if err != nil {
		t.Fatalf("Failed to delete boot: %v", err)
	}
	t.Log("DeleteBoot: Deleted!")
}
