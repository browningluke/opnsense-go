package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewClientEnablesHTTP2(t *testing.T) {
	client := NewClient(Options{})
	transport, ok := client.client.HTTPClient.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("transport type = %T, want *http.Transport", client.client.HTTPClient.Transport)
	}
	if !transport.ForceAttemptHTTP2 {
		t.Fatal("ForceAttemptHTTP2 is disabled")
	}
}

func TestDoRequestDoesNotLogPayloads(t *testing.T) {
	const requestSecret = "request-secret"
	const responseSecret = "response-secret"
	const sessionSecret = "session-secret"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.SetCookie(w, &http.Cookie{Name: "PHPSESSID", Value: sessionSecret})
		_ = json.NewEncoder(w).Encode(map[string]string{"secret": responseSecret})
	}))
	defer server.Close()

	var logs bytes.Buffer
	client := NewClient(Options{
		Uri:       server.URL,
		APIKey:    "api-key",
		APISecret: "api-secret",
		Logger:    log.New(&logs, "", 0),
	})

	var response map[string]string
	if err := client.doRequest(context.Background(), http.MethodPost, "/test", map[string]string{"secret": requestSecret}, &response); err != nil {
		t.Fatalf("doRequest() error = %v", err)
	}
	if response["secret"] != responseSecret {
		t.Fatalf("response secret = %q, want %q", response["secret"], responseSecret)
	}

	for _, secret := range []string{requestSecret, responseSecret, sessionSecret, "api-key", "api-secret"} {
		if strings.Contains(logs.String(), secret) {
			t.Errorf("logs contain sensitive value %q: %s", secret, logs.String())
		}
	}
}
