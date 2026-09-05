package zenarmor

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestPolicies(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet || request.URL.Path != "/api/zenarmor/policy" {
			t.Fatalf("unexpected request: %s %s", request.Method, request.URL.Path)
		}
		_ = json.NewEncoder(response).Encode(PoliciesResponse{Policies: []Policy{{ID: 42, LocalID: 42, Name: "test", IsActive: true}}})
	}))
	defer server.Close()

	controller := &Controller{Api: api.NewClient(api.Options{Uri: server.URL, APIKey: "key", APISecret: "secret", Logger: log.New(io.Discard, "", 0)})}
	policies, err := controller.Policies(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(policies) != 1 || policies[0].ID != 42 || policies[0].Name != "test" || !policies[0].IsActive {
		t.Fatalf("unexpected policies: %#v", policies)
	}
}
