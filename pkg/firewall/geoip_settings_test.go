package firewall

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestGeoIPSettings(t *testing.T) {
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

	// Step 1: Get current settings (equivalent to import step).
	respGet, err := controller.GeoIPSettingsGet(ctx)
	if err != nil {
		t.Fatalf("Failed to get GeoIP settings: %v", err)
	}
	t.Logf("GeoIPSettingsGet (initial): %+v", respGet)

	origUrl := respGet.Alias.GeoIP.Url

	// Step 2: Set a new URL.
	updatedUrl := "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-Country-CSV&license_key=test&suffix=zip"
	_, err = controller.GeoIPSettingsSet(ctx, &GeoIPSettingsSetParams{
		GeoIP: GeoIPSettingsUrl{Url: updatedUrl},
	})
	if err != nil {
		t.Fatalf("Failed to set GeoIP settings: %v", err)
	}

	// Step 3: Read back and verify the change was applied.
	respGet, err = controller.GeoIPSettingsGet(ctx)
	if err != nil {
		t.Fatalf("Failed to get GeoIP settings after set: %v", err)
	}
	if respGet.Alias.GeoIP.Url != updatedUrl {
		t.Fatalf("url not updated; got %q, want %q", respGet.Alias.GeoIP.Url, updatedUrl)
	}

	// Step 4: Restore original value.
	_, err = controller.GeoIPSettingsSet(ctx, &GeoIPSettingsSetParams{
		GeoIP: GeoIPSettingsUrl{Url: origUrl},
	})
	if err != nil {
		t.Fatalf("Failed to restore GeoIP settings: %v", err)
	}

	// Step 5: Verify restore.
	respGet, err = controller.GeoIPSettingsGet(ctx)
	if err != nil {
		t.Fatalf("Failed to get GeoIP settings after restore: %v", err)
	}
	if respGet.Alias.GeoIP.Url != origUrl {
		t.Fatalf("url not restored; got %q, want %q", respGet.Alias.GeoIP.Url, origUrl)
	}
}
