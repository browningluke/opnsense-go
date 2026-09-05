package core

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestFirmwareInfo(t *testing.T) {
	if os.Getenv("OPNSENSE_URI") == "" || os.Getenv("OPNSENSE_API_KEY") == "" || os.Getenv("OPNSENSE_API_SECRET") == "" {
		t.Skip("OPNSENSE API environment variables are required")
	}
	controller := &Controller{Api: api.NewClient(api.Options{
		Uri:           os.Getenv("OPNSENSE_URI"),
		APIKey:        os.Getenv("OPNSENSE_API_KEY"),
		APISecret:     os.Getenv("OPNSENSE_API_SECRET"),
		AllowInsecure: os.Getenv("OPNSENSE_ALLOW_INSECURE") == "true",
	})}

	info, err := controller.FirmwareInfo(context.Background())
	if err != nil {
		t.Fatalf("FirmwareInfo failed: %v", err)
	}
	if len(info.Plugin) == 0 {
		t.Fatal("FirmwareInfo returned no plugin inventory")
	}
}
