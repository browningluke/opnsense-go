package openvpn

import (
	"context"
	"testing"
)

func TestServiceReconfigure(t *testing.T) {
	c := newTestController(t)
	ctx := context.Background()

	res, err := c.ServiceReconfigure(ctx)
	if err != nil {
		t.Fatalf("ServiceReconfigure failed: %v", err)
	}
	if res.Result != "ok" {
		t.Fatalf("expected result ok; got %q", res.Result)
	}
}

func TestServiceGenKey(t *testing.T) {
	c := newTestController(t)
	ctx := context.Background()

	// Default form: no ?type= query parameter (returns a "secret" static key).
	res, err := c.ServiceGenKey(ctx, nil)
	if err != nil {
		t.Fatalf("ServiceGenKey(nil) failed: %v", err)
	}
	if res.Result != "ok" {
		t.Fatalf("ServiceGenKey(nil): expected result ok; got %q", res.Result)
	}
	if res.Key == "" {
		t.Fatal("ServiceGenKey(nil) returned empty key")
	}

	// Each hyphenated key type — the case that motivated the query-parameter
	// support: passed as ?type=... rather than in the path.
	for _, kt := range []string{"tls-auth", "tls-crypt", "tls-crypt-v2-server"} {
		kt := kt
		res, err := c.ServiceGenKey(ctx, &kt)
		if err != nil {
			t.Fatalf("ServiceGenKey(%s) failed: %v", kt, err)
		}
		if res.Result != "ok" {
			t.Fatalf("ServiceGenKey(%s): expected result ok; got %q", kt, res.Result)
		}
		if res.Key == "" {
			t.Fatalf("ServiceGenKey(%s) returned empty key", kt)
		}
	}
}
