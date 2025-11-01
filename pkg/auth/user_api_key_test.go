package auth

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestUserApiKeys(t *testing.T) {
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
		Name:     "integration",
		Disabled: "0",
		Password: "Integration123",
		Fullname: "Integration Test",
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

	respAdd, err := controller.UserAddApiKey(ctx, respUser.Name)
	if err != nil {
		t.Fatalf("Failed to add api key for user: %v", err)
	}
	t.Logf("AddApiKey: %+v", respAdd)

	respGet, err := controller.UserGetAllApiKeys(ctx)
	if err != nil {
		t.Fatalf("Failed to get all user api keys: %v", err)
	}
	t.Logf("GetAllApiKeysResult: %+v", respGet)

	var id string
	for _, v := range respGet.Rows {
		if v.Username == respUser.Name {
			id = v.Id
			break
		}
	}
	respDel, err := controller.UserDeleteApiKey(ctx, id)
	if err != nil {
		t.Fatalf("Failed to delete api key for user: %v", err)
	}
	t.Logf("DeleteApiKey: %+v", respDel)

	err = controller.DeleteUser(ctx, key)
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}
	t.Logf("Deleted user with key: %s", key)
}
