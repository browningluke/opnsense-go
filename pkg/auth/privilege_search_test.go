package auth

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestPrivilegeSearch(t *testing.T) {
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

	resp, err := controller.PrivilegeGetAll(ctx)
	if err != nil {
		t.Fatalf("Failed to search privileges: %v", err)
	}
	t.Logf("PrivilegeSearchResult: %+v", resp)
}
