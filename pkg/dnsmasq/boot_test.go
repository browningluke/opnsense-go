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
		// Interface: api.SelectedMap("lan"),
		// Tag:           api.SelectedMapList([]string{"b4b79319-2e09-47ae-85c2-687eb4b6e7ee", "d594fa8a-1a76-44e7-afab-adb6c5bdb69e"}),
		Filename:      "test",
		Servername:    "test-servername",
		ServerAddress: "test-serveraddress",
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

	// boot.Interface = api.SelectedMap("wan")
	// boot.Tag = api.SelectedMapList([]string{"b4b79319-2e09-47ae-85c2-687eb4b6e7ee"})
	boot.Filename = "test-updated"
	boot.Servername = "test-servername-updated"
	boot.ServerAddress = "test-serveraddress-updated"
	boot.Description = "test-description-updated"
	err = controller.UpdateBoot(ctx, respAdd, boot)
	if err != nil {
		t.Fatalf("Failed to update boot: %v", err)
	}
	t.Logf("UpdateBoot: %+v", boot)

	respSearch, err := controller.SearchBoot(ctx, "-1")
	if err != nil {
		t.Fatalf("Failed to search Boots: %v", err)
	}
	t.Logf("SearchBoot: %+v", respSearch)
	noRowFound := true
	lastId := ""
	for _, v := range respSearch.Rows {
		lastId = v.Id
		if v.Id != respAdd {
			continue
		}
		noRowFound = false
		if v.InterfaceId != boot.Interface.String() {
			t.Fatalf("Interface not updated; Got: %s Expected: %s", v.Interface, boot.Interface.String())
		}
		if v.Filename != boot.Filename {
			t.Fatalf("Port not updated; Got: %s Expected: %s", v.Filename, boot.Filename)
		}
		if v.Servername != boot.Servername {
			t.Fatalf("Servername not updated; Got: %s Expected: %s", v.Servername, boot.Servername)
		}
		if v.ServerAddress != boot.ServerAddress {
			t.Fatalf("ServerAddress not updated; Got: %s Expected: %s", v.ServerAddress, boot.ServerAddress)
		}
		if v.Description != boot.Description {
			t.Fatalf("Description not updated; Got: %s Expected: %s", v.Description, boot.Description)
		}
	}
	if noRowFound {
		t.Fatalf("Row not found that was added; Got: %s Expected: %s", lastId, respAdd)
	}

	err = controller.DeleteBoot(ctx, respAdd)
	if err != nil {
		t.Fatalf("Failed to delete boot: %v", err)
	}
	t.Log("DeleteBoot: Deleted!")
}
