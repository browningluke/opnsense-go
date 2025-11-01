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
		Interface: api.SelectedMap(""),
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

	respSearch, err := controller.SearchOption(ctx, "-1")
	if err != nil {
		t.Fatalf("Failed to search Options: %v", err)
	}
	t.Logf("SearchOption: %+v", respSearch)
	noRowFound := true
	lastId := ""
	for _, v := range respSearch.Rows {
		lastId = v.Id
		if v.Id != respAdd {
			continue
		}
		noRowFound = false
		if v.OptionV4Id != option.OptionV4.String() {
			t.Fatalf("OptionV4 not updated; Got: %s Expected: %s", v.OptionV4Id, option.OptionV4.String())
		}
		if v.OptionV6Id != option.OptionV6.String() {
			t.Fatalf("OptionV6 not updated; Got: %s Expected: %s", v.OptionV6Id, option.OptionV4.String())
		}
		if v.Value != option.Value {
			t.Fatalf("Value not updated; Got: %s Expected: %s", v.Value, option.Value)
		}
		if v.Description != option.Description {
			t.Fatalf("Description not updated; Got: %s Expected: %s", v.Description, option.Description)
		}
	}
	if noRowFound {
		t.Fatalf("Row not found that was added; Got: %s Expected: %s", lastId, respAdd)
	}

	err = controller.DeleteOption(ctx, respAdd)
	if err != nil {
		t.Fatalf("Failed to delete Option: %v", err)
	}
	t.Log("DeleteOption: Deleted!")
}
