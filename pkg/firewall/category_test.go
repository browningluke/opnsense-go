package firewall

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestCategory(t *testing.T) {
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

	category := &Category{
		Automatic: "0",
		Name:      "test-category",
		Color:     "FF0000",
	}

	key, err := controller.AddCategory(ctx, category)
	if err != nil {
		t.Fatalf("Failed to add category: %v", err)
	}
	t.Logf("Added category with key: %s", key)

	retrievedCategory, err := controller.GetCategory(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get category: %v", err)
	}
	t.Logf("Retrieved category: %+v", retrievedCategory)
	if retrievedCategory.Name != category.Name {
		t.Fatalf("Retrieved category name does not match: got %s, want %s", retrievedCategory.Name, category.Name)
	}
	if retrievedCategory.Color != category.Color {
		t.Fatalf("Retrieved category color does not match: got %s, want %s", retrievedCategory.Color, category.Color)
	}
	if retrievedCategory.Automatic != category.Automatic {
		t.Fatalf("Retrieved category automatic does not match: got %s, want %s", retrievedCategory.Automatic, category.Automatic)
	}

	category.Name = "test-category-updated"
	category.Color = "00FF00"
	err = controller.UpdateCategory(ctx, key, category)
	if err != nil {
		t.Fatalf("Failed to update category: %v", err)
	}
	t.Logf("Updated category with key: %s", key)

	retrievedCategory, err = controller.GetCategory(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get updated category: %v", err)
	}
	if retrievedCategory.Name != "test-category-updated" {
		t.Fatalf("Retrieved category name does not match updated name: got %s, want %s", retrievedCategory.Name, "test-category-updated")
	}
	if retrievedCategory.Color != "00FF00" {
		t.Fatalf("Retrieved category color does not match updated color: got %s, want %s", retrievedCategory.Color, "00FF00")
	}

	err = controller.DeleteCategory(ctx, key)
	if err != nil {
		t.Fatalf("Failed to delete category: %v", err)
	}
	t.Logf("Deleted category with key: %s", key)
}
