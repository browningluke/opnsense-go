package auth

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestPrivileges(t *testing.T) {
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

	respItem, err := controller.PrivilegeGetItem(ctx, "page-all")
	if err != nil {
		t.Fatalf("Failed to get privilege item: %+v", err)
	}
	t.Logf("Result: %+v", respItem)

	groupId := respItem.Item.Groups.String()
	item := PrivilegeSetItem{
		Groups: groupId,
		Users:  "",
	}
	respSet, err := controller.PrivilegeSetItem(ctx, "page-diagnostics-arptable", &item)
	if err != nil {
		t.Fatalf("Failed to set privilege item: %+v", err)
	}
	if respSet.Result != "saved" {
		t.Fatalf("Failed to set privilege item: %+v", item)
	}
	t.Logf("PrivilegeSetItem: %+v", respSet)

	respPagePriv, err := controller.PrivilegeGetItem(ctx, "page-diagnostics-arptable")
	if err != nil {
		t.Fatalf("Failed to get privilege item: %+v", err)
	}
	if respPagePriv.Item.Groups.String() != groupId {
		t.Fatalf("Failed to set privilege group on item: %+v", respPagePriv)
	}
	t.Logf("PrivilegeGetItem: %+v", respPagePriv)

	item.Groups = ""
	respReset, err := controller.PrivilegeSetItem(ctx, "page-diagnostics-arptable", &item)
	if err != nil {
		t.Fatalf("Failed to set privilege item: %+v", err)
	}
	if respReset.Result != "saved" {
		t.Fatalf("Failed to reset privilege item: %+v", respReset)
	}
	t.Logf("PrivilegeSetItem: %+v", respReset)
}
