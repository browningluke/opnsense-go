package core

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/browningluke/opnsense-go/pkg/api"
)

// testFirmwarePkg is a small, fast-installing community plugin with no reboot
// requirement. Chosen specifically so the test can install→verify→remove
// without leaving the firewall in an unexpected state.
const testFirmwarePkg = "os-acme-client"

// pollInterval/pollTimeout are sized for a real OPNsense install of a small
// plugin: typical wall-clock 20-60s on a low-end box, dominated by fetch.
// A 5-minute ceiling absorbs slow mirrors without making the test loop chatty.
const (
	pollInterval = 3 * time.Second
	pollTimeout  = 5 * time.Minute
)

func newFirmwareController(t *testing.T) *Controller {
	t.Helper()
	uri := os.Getenv("OPNSENSE_URI")
	if uri == "" {
		t.Skip("OPNSENSE_URI not set; skipping firmware integration test")
	}
	return &Controller{
		Api: api.NewClient(api.Options{
			Uri:           uri,
			APIKey:        os.Getenv("OPNSENSE_API_KEY"),
			APISecret:     os.Getenv("OPNSENSE_API_SECRET"),
			AllowInsecure: true,
			MaxBackoff:    30,
			MinBackoff:    1,
			MaxRetries:    4,
		}),
	}
}

// waitForFirmwareIdle blocks until /firmware/running reports "ready" or the
// shared pollTimeout elapses. The OPNsense firmware backend has a single
// global job slot; submitting a new install while another job is still
// finalising will be silently queued or rejected. Always block on this before
// any mutating call. Returns an error rather than calling t.Fatal so cleanup
// callers (where t.Fatal is unsafe) can use the same function.
func waitForFirmwareIdle(ctx context.Context, c *Controller) error {
	deadline, cancel := context.WithTimeout(ctx, pollTimeout)
	defer cancel()
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	for {
		s, err := c.FirmwareRunning(deadline)
		if err != nil {
			return err
		}
		if s.Status == "ready" {
			return nil
		}
		select {
		case <-deadline.Done():
			return deadline.Err()
		case <-ticker.C:
		}
	}
}

// waitForFirmwareJobDone blocks until /firmware/upgradestatus reports a
// terminal status ("done"/"reboot"/"error") or the shared pollTimeout
// elapses. Returns the final status struct so the caller can inspect the log.
// Returns an error rather than calling t.Fatal so cleanup callers can reuse
// it.
func waitForFirmwareJobDone(ctx context.Context, c *Controller) (*FirmwareUpgradeStatus, error) {
	deadline, cancel := context.WithTimeout(ctx, pollTimeout)
	defer cancel()
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	for {
		s, err := c.FirmwareUpgradeStatus(deadline)
		if err != nil {
			return nil, err
		}
		switch s.Status {
		case "done", "reboot":
			return s, nil
		case "error":
			return s, &firmwareJobError{log: s.Log}
		}
		select {
		case <-deadline.Done():
			return nil, deadline.Err()
		case <-ticker.C:
		}
	}
}

type firmwareJobError struct{ log string }

func (e *firmwareJobError) Error() string {
	// Trim the log to the last few lines so test output isn't flooded.
	lines := strings.Split(strings.TrimSpace(e.log), "\n")
	if len(lines) > 10 {
		lines = lines[len(lines)-10:]
	}
	return "firmware job reported error; last lines: " + strings.Join(lines, " | ")
}

// findPlugin returns the plugin entry for name, or nil if not present in info.
func findPlugin(info *FirmwareInfo, name string) *FirmwarePlugin {
	for i := range info.Plugins {
		if info.Plugins[i].Name == name {
			return &info.Plugins[i]
		}
	}
	return nil
}

// TestFirmware exercises the plugin lifecycle: install → verify → lock →
// unlock → remove. Requires a live OPNsense reachable via OPNSENSE_URI. The
// chosen test plugin (os-acme-client) is small, fast, and reboot-free.
//
// The test refuses to run if testFirmwarePkg is already installed at start,
// to avoid masking an Install failure and to keep the test fully reversible.
func TestFirmware(t *testing.T) {
	c := newFirmwareController(t)
	ctx := context.Background()

	// Pre-check: backend idle.
	if err := waitForFirmwareIdle(ctx, c); err != nil {
		t.Fatalf("waiting for firmware idle: %v", err)
	}

	// Pre-check: package is known to the remote and not currently installed.
	info, err := c.FirmwareInfo(ctx)
	if err != nil {
		t.Fatalf("FirmwareInfo: %v", err)
	}
	p := findPlugin(info, testFirmwarePkg)
	if p == nil {
		t.Fatalf("plugin %s not advertised by mirror; refusing to test", testFirmwarePkg)
	}
	if p.Installed == "1" {
		t.Fatalf("plugin %s already installed; refusing to test (would leak state on cleanup)", testFirmwarePkg)
	}

	// Always attempt to leave the firewall how we found it, even if a step
	// below fails partway through. Cleanup is best-effort: a remove of an
	// uninstalled package is harmless to opnsense.
	t.Cleanup(func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), pollTimeout)
		defer cancel()
		if err := waitForFirmwareIdle(cleanupCtx, c); err != nil {
			t.Logf("cleanup: waiting for firmware idle: %v", err)
			return
		}
		if _, err := c.FirmwareRemove(cleanupCtx, testFirmwarePkg); err != nil {
			t.Logf("cleanup: FirmwareRemove(%s): %v", testFirmwarePkg, err)
			return
		}
		if _, err := waitForFirmwareJobDone(cleanupCtx, c); err != nil {
			t.Logf("cleanup: waiting for remove to finish: %v", err)
		}
	})

	// Install.
	res, err := c.FirmwareInstall(ctx, testFirmwarePkg)
	if err != nil {
		t.Fatalf("FirmwareInstall(%s): %v", testFirmwarePkg, err)
	}
	if res.Status != "ok" {
		t.Fatalf("FirmwareInstall status: got %q, want %q", res.Status, "ok")
	}
	t.Logf("install accepted; msg_uuid=%s", res.MsgUUID)
	final, err := waitForFirmwareJobDone(ctx, c)
	if err != nil {
		t.Fatalf("waiting for install: %v", err)
	}
	t.Logf("install finished; status=%s", final.Status)

	// Verify installed.
	info, err = c.FirmwareInfo(ctx)
	if err != nil {
		t.Fatalf("FirmwareInfo (post-install): %v", err)
	}
	p = findPlugin(info, testFirmwarePkg)
	if p == nil || p.Installed != "1" {
		t.Fatalf("post-install: plugin %s not reported installed (entry=%+v)", testFirmwarePkg, p)
	}

	// Lock.
	if err := waitForFirmwareIdle(ctx, c); err != nil {
		t.Fatalf("waiting for firmware idle pre-lock: %v", err)
	}
	if res, err = c.FirmwareLock(ctx, testFirmwarePkg); err != nil {
		t.Fatalf("FirmwareLock(%s): %v", testFirmwarePkg, err)
	}
	if res.Status != "ok" {
		t.Fatalf("FirmwareLock status: got %q, want %q", res.Status, "ok")
	}
	if _, err := waitForFirmwareJobDone(ctx, c); err != nil {
		t.Fatalf("waiting for lock: %v", err)
	}

	info, err = c.FirmwareInfo(ctx)
	if err != nil {
		t.Fatalf("FirmwareInfo (post-lock): %v", err)
	}
	p = findPlugin(info, testFirmwarePkg)
	if p == nil {
		t.Fatalf("post-lock: plugin %s vanished from info", testFirmwarePkg)
	}
	// OPNsense reports locked as a string-ish field; treat any non-"N/A"/
	// non-"0" as locked since the exact value is implementation-defined.
	if p.Locked == "N/A" || p.Locked == "0" || p.Locked == "" {
		t.Fatalf("post-lock: expected locked, got Locked=%q", p.Locked)
	}

	// Unlock.
	if err := waitForFirmwareIdle(ctx, c); err != nil {
		t.Fatalf("waiting for firmware idle pre-unlock: %v", err)
	}
	if res, err = c.FirmwareUnlock(ctx, testFirmwarePkg); err != nil {
		t.Fatalf("FirmwareUnlock(%s): %v", testFirmwarePkg, err)
	}
	if res.Status != "ok" {
		t.Fatalf("FirmwareUnlock status: got %q, want %q", res.Status, "ok")
	}
	if _, err := waitForFirmwareJobDone(ctx, c); err != nil {
		t.Fatalf("waiting for unlock: %v", err)
	}

	info, err = c.FirmwareInfo(ctx)
	if err != nil {
		t.Fatalf("FirmwareInfo (post-unlock): %v", err)
	}
	p = findPlugin(info, testFirmwarePkg)
	if p == nil {
		t.Fatalf("post-unlock: plugin %s vanished from info", testFirmwarePkg)
	}

	// Remove (the deferred cleanup would do this too, but we want to verify
	// the result rather than swallow errors).
	if err := waitForFirmwareIdle(ctx, c); err != nil {
		t.Fatalf("waiting for firmware idle pre-remove: %v", err)
	}
	if res, err = c.FirmwareRemove(ctx, testFirmwarePkg); err != nil {
		t.Fatalf("FirmwareRemove(%s): %v", testFirmwarePkg, err)
	}
	if res.Status != "ok" {
		t.Fatalf("FirmwareRemove status: got %q, want %q", res.Status, "ok")
	}
	if _, err := waitForFirmwareJobDone(ctx, c); err != nil {
		t.Fatalf("waiting for remove: %v", err)
	}

	info, err = c.FirmwareInfo(ctx)
	if err != nil {
		t.Fatalf("FirmwareInfo (post-remove): %v", err)
	}
	p = findPlugin(info, testFirmwarePkg)
	if p != nil && p.Installed == "1" {
		t.Fatalf("post-remove: plugin %s still reported installed", testFirmwarePkg)
	}
}

// TestFirmwareInfo is a non-mutating sanity check that any caller with a live
// OPNsense can run safely; verifies the response unmarshals and contains the
// expected top-level shape (at least one package, at least one plugin).
func TestFirmwareInfo(t *testing.T) {
	c := newFirmwareController(t)
	ctx := context.Background()

	info, err := c.FirmwareInfo(ctx)
	if err != nil {
		t.Fatalf("FirmwareInfo: %v", err)
	}
	if info.ProductID == "" {
		t.Errorf("ProductID empty; expected non-empty")
	}
	if info.ProductVersion == "" {
		t.Errorf("ProductVersion empty; expected non-empty")
	}
	if len(info.Packages) == 0 {
		t.Errorf("Packages empty; OPNsense install should always have packages")
	}
	// Plugins may legitimately be empty on a minimal install, so don't
	// assert on length — but if present, sanity-check the first entry.
	if len(info.Plugins) > 0 && info.Plugins[0].Name == "" {
		t.Errorf("first plugin has empty Name; struct tags likely wrong")
	}
}

// TestFirmwareRunning is a non-mutating sanity check that running status
// decodes and reports the documented idle value.
func TestFirmwareRunning(t *testing.T) {
	c := newFirmwareController(t)
	ctx := context.Background()

	s, err := c.FirmwareRunning(ctx)
	if err != nil {
		t.Fatalf("FirmwareRunning: %v", err)
	}
	if s.Status == "" {
		t.Errorf("Status empty; expected at least 'ready'")
	}
}
