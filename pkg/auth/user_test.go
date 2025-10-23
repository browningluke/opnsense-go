package auth

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestUser(t *testing.T) {
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

	user := &User{
		Name:              "test",
		Disabled:          "0",
		Password:          "test123",
		ScrambledPassword: "0",
		Comment:           "Test comment",
		Email:             "test@test.nl",
		Fullname:          "Test description full name",
	}

	key, err := controller.AddUser(ctx, user)
	if err != nil {
		t.Fatalf("Failed to add user: %v", err)
	}
	t.Logf("Added user with key: %s", key)

	respUser, err := controller.GetUser(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}
	t.Logf("Retrieved user: %+v", respUser)
	if respUser.Name != user.Name {
		t.Fatalf("Retrieved user name does not match: got %s, want %s", respUser.Name, user.Name)
	}

	user.Name = "test-username"
	user.Disabled = "1"
	user.Shell = api.SelectedMap("/bin/csh")
	user.Priviledge = api.SelectedMapList([]string{"page-diagnostics-authentication"})
	user.GroupMemberships = api.SelectedMapList([]string{"1999"})
	user.OtpSeed = "Z2A3Y4EERVKCG24W6Q3VQUUKTK7HZFVM"
	user.AuthorizedKeys = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAINYqRfv6UauTujjdiwZAUALv38Z0OXII20h9q6KdvbyZ mike@mike-desktop-pc"

	err = controller.UpdateUser(ctx, key, user)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}
	t.Logf("Updated user with key: %s", key)

	respUser, err = controller.GetUser(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}
	t.Logf("Retrieved user: %+v", respUser)

	err = controller.DeleteUser(ctx, key)
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}
	t.Logf("Deleted user with key: %s", key)
}
