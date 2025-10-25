package dnsmasq

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestTag(t *testing.T) {
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

	tag := &Tag{
		Tag: "2",
	}

	respAdd, err := controller.AddTag(ctx, tag)
	if err != nil {
		t.Fatalf("Failed to add tag: %v", err)
	}
	t.Logf("AddTag: %+v", respAdd)

	respGet, err := controller.GetTag(ctx, respAdd)
	if err != nil {
		t.Fatalf("Failed to get tag: %v", err)
	}
	t.Logf("GetTag: %+v", respGet)

	tag.Tag = "3"
	err = controller.UpdateTag(ctx, respAdd, tag)
	if err != nil {
		t.Fatalf("Failed to update tag: %v", err)
	}
	t.Logf("UpdateTag: %+v", tag)

	respGet, err = controller.GetTag(ctx, respAdd)
	if err != nil {
		t.Fatalf("Failed to get tag: %v", err)
	}
	t.Logf("GetTag: %+v", respGet)

	err = controller.DeleteTag(ctx, respAdd)
	if err != nil {
		t.Fatalf("Failed to delete tag: %v", err)
	}
	t.Log("DeleteTag: Deleted!")
}
