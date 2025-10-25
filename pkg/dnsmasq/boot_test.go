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

	boot := &Boot{}

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
