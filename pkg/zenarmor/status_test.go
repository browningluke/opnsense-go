package zenarmor

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestStatus(t *testing.T) {
	t.Parallel()
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/zenarmor/status" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"agent_installed":true,"agent_enabled":true,"agent_status":true,"agent_version":{"version":"2.6.1","lastUpdate":1},"engine":{"version":"1.18","lastUpdate":2},"eastpect":{"status":true},"database_version":{"version":"2026.08","lastUpdate":3},"database":{"info":{"version":"7.0"},"status":true},"license":"Free","cloud_threat_intel":false,"deployment_mode":"local"}`))
	}))
	defer server.Close()

	client := api.NewClient(api.Options{Uri: server.URL, APIKey: "key", APISecret: "secret", AllowInsecure: true, Logger: log.New(io.Discard, "", 0)})
	status, err := (&Controller{Api: client}).Status(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !status.AgentStatus || !status.Eastpect.Status || status.Engine.Version != "1.18" || status.DatabaseVersion.Version != "2026.08" || status.License != "Free" {
		t.Fatalf("unexpected status: %#v", status)
	}
}
