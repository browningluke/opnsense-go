package dnsmasq

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestOption(t *testing.T) {
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

	option := &Option{
		Type:     api.SelectedMap("set"),
		OptionV4: api.SelectedMap("1"),
		// OptionV6: api.SelectedMap("23"),
		Interface: api.SelectedMap("lan"),
		// Tag:         api.SelectedMapList([]string{"3dff7e13-68e6-4d37-a9ac-35944c1cc116"}),
		Value:       "255.255.255.0",
		Description: "test subnetmask",
	}

	respAdd, err := controller.AddOption(ctx, option)
	if err != nil {
		t.Fatalf("Failed to add Option: %v", err)
	}
	t.Logf("AddOption: %+v", respAdd)

	respGet, err := controller.GetOption(ctx, respAdd)
	if err != nil {
		t.Fatalf("Failed to get Option: %v", err)
	}
	t.Logf("GetOption: %+v", respGet)

	// option.Type = api.SelectedMap("match")
	// option.OptionV4 = api.SelectedMap("3")
	option.OptionV4 = api.SelectedMap("")
	option.OptionV6 = api.SelectedMap("23")
	// option.Tag = nil
	// option.TagSet = api.SelectedMap("763bd6ed-e6bc-4c3c-aef7-6fef954179a5")
	option.Value = "255.255.255.250"
	option.Description = "test-subnetmask-updated"
	err = controller.UpdateOption(ctx, respAdd, option)
	if err != nil {
		t.Fatalf("Failed to update Option: %v", err)
	}
	t.Logf("UpdateOption: %+v", option)

	respGet, err = controller.GetOption(ctx, respAdd)
	if err != nil {
		t.Fatalf("Failed to get Option: %v", err)
	}
	t.Logf("GetOption: %+v", respGet)

	err = controller.DeleteOption(ctx, respAdd)
	if err != nil {
		t.Fatalf("Failed to delete Option: %v", err)
	}
	t.Log("DeleteOption: Deleted!")
}
