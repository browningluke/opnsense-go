package dnsmasq

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/firewall"
)

func TestIntegration(t *testing.T) {
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

	tagOne := &Tag{
		Tag: "1",
	}
	tagTwo := &Tag{
		Tag: "2",
	}

	tagOneId, err := controller.AddTag(ctx, tagOne)
	if err != nil || tagOneId == "" {
		t.Fatalf("Error inserting: %+v; %s", tagOne, err)
	}

	tagTwoId, err := controller.AddTag(ctx, tagTwo)
	if err != nil || tagTwoId == "" {
		t.Fatalf("Error inserting: %+v; %s", tagTwo, err)
	}

	optionOne := &Option{
		Type:        api.SelectedMap("set"),
		OptionV4:    api.SelectedMap("1"), //netmask
		Interface:   api.SelectedMap(""),
		TypeSetTags: api.SelectedMapList([]string{tagOneId}),
	}

	optionTwo := &Option{
		Type:         api.SelectedMap("match"),
		OptionV6:     api.SelectedMap("32"), //information-refresh-time
		TypeMatchTag: api.SelectedMap(tagTwoId),
	}

	optionOneId, err := controller.AddOption(ctx, optionOne)
	if err != nil || optionOneId == "" {
		t.Fatalf("Error inserting: %+v; %s", optionOne, err)
	}
	optionTwoId, err := controller.AddOption(ctx, optionTwo)
	if err != nil || optionTwoId == "" {
		t.Fatalf("Error inserting: %+v; %s", optionTwo, err)
	}

	optionOne.TypeSetTags = append(optionOne.TypeSetTags, tagTwoId)
	err = controller.UpdateOption(ctx, optionOneId, optionOne)
	if err != nil {
		t.Fatalf("Error updating: %+v; %s", optionOne, err)
	}

	optionTwo.TypeMatchTag = api.SelectedMap(tagOneId)
	controller.UpdateOption(ctx, optionTwoId, optionTwo)
	if err != nil {
		t.Fatalf("Error updating: %+v; %s", optionTwo, err)
	}

	bootOne := &Boot{
		Tag:      api.SelectedMapList([]string{tagOneId}),
		Filename: "test-filename",
	}

	bootOneId, err := controller.AddBoot(ctx, bootOne)
	if err != nil || bootOneId == "" {
		t.Fatalf("Error inserting: %+v; %s", bootOne, err)
	}

	bootOne.Tag = append(bootOne.Tag, tagTwoId)
	err = controller.UpdateBoot(ctx, bootOneId, bootOne)
	if err != nil {
		t.Fatalf("Error updating: %+v; %s", bootOne, err)
	}

	rangeOne := &Range{
		StartAddress: "192.168.100.200",
		EndAddress:   "192.168.100.250",
		Tag:          api.SelectedMap(tagOneId),
		DomainType:   api.SelectedMap("range"),
	}

	rangeOneId, err := controller.AddRange(ctx, rangeOne)
	if err != nil || rangeOneId == "" {
		t.Fatalf("Error inserting: %+v; %s", rangeOne, err)
	}

	rangeOne.Tag = api.SelectedMap(tagTwoId)
	err = controller.UpdateRange(ctx, rangeOneId, rangeOne)
	if err != nil {
		t.Fatalf("Error updating: %+v; %s", rangeOne, err)
	}

	firewallAliasOne := &firewall.Alias{
		Enabled: "1",
		Name:    "test_firewall_alias_one",
		Type:    api.SelectedMap("external"),
	}

	firewallAliasTwo := &firewall.Alias{
		Enabled: "1",
		Name:    "test_firewall_alias_two",
		Type:    api.SelectedMap("external"),
	}

	fwController := firewall.Controller{
		Api: api_client,
	}
	firewallAliasOneId, err := fwController.AddAlias(ctx, firewallAliasOne)
	if err != nil || firewallAliasOneId == "" {
		t.Fatalf("Error inserting: %+v; %s", firewallAliasOne, err)
	}

	firewallAliasTwoId, err := fwController.AddAlias(ctx, firewallAliasTwo)
	if err != nil || firewallAliasTwoId == "" {
		t.Fatalf("Error inserting: %+v; %s", firewallAliasTwo, err)
	}

	domainOne := &Domain{
		Sequence:      "1",
		Domain:        "test-domain",
		FirewallAlias: api.SelectedMap(firewallAliasOneId),
	}

	domainOneId, err := controller.AddDomain(ctx, domainOne)
	if err != nil || domainOneId == "" {
		t.Fatalf("Error inserting: %+v; %s", domainOne, err)
	}

	domainOne.FirewallAlias = api.SelectedMap(firewallAliasTwoId)
	err = controller.UpdateDomain(ctx, domainOneId, domainOne)
	if err != nil {
		t.Fatalf("Error updating: %+v; %s", domainOne, err)
	}

	hostOne := &Host{
		Hostname:    "test-hostname",
		IpAddresses: api.SelectedMapList([]string{"192.168.100.100", "fe80::0001"}),
		Tag:         api.SelectedMap(tagOneId),
	}

	hostOneId, err := controller.AddHost(ctx, hostOne)
	if err != nil || hostOneId == "" {
		t.Fatalf("Error inserting: %+v; %s", hostOne, err)
	}

	hostOne.Tag = api.SelectedMap(tagTwoId)
	err = controller.UpdateHost(ctx, hostOneId, hostOne)
	if err != nil {
		t.Fatalf("Error updating: %+v; %s", hostOne, err)
	}

	//Cleanup
	err = controller.DeleteHost(ctx, hostOneId)
	if err != nil {
		t.Fatalf("Error deleting: %+v; %s", hostOne, err)
	}

	err = controller.DeleteDomain(ctx, domainOneId)
	if err != nil {
		t.Fatalf("Error deleting: %+v; %s", domainOne, err)
	}

	err = fwController.DeleteAlias(ctx, firewallAliasOneId)
	if err != nil {
		t.Fatalf("Error deleting: %+v; %s", firewallAliasOne, err)
	}
	err = fwController.DeleteAlias(ctx, firewallAliasTwoId)
	if err != nil {
		t.Fatalf("Error deleting: %+v; %s", firewallAliasTwo, err)
	}

	err = controller.DeleteRange(ctx, rangeOneId)
	if err != nil {
		t.Fatalf("Error deleting: %+v; %s", rangeOne, err)
	}

	err = controller.DeleteBoot(ctx, bootOneId)
	if err != nil {
		t.Fatalf("Error deleting: %+v; %s", bootOne, err)
	}

	err = controller.DeleteOption(ctx, optionOneId)
	if err != nil {
		t.Fatalf("Error deleting: %+v; %s", optionOne, err)
	}
	err = controller.DeleteOption(ctx, optionTwoId)
	if err != nil {
		t.Fatalf("Error deleting: %+v; %s", optionTwo, err)
	}

	err = controller.DeleteTag(ctx, tagOneId)
	if err != nil {
		t.Fatalf("Error deleting: %+v; %s", tagOne, err)
	}
	err = controller.DeleteTag(ctx, tagTwoId)
	if err != nil {
		t.Fatalf("Error deleting: %+v; %s", tagTwo, err)
	}
}
