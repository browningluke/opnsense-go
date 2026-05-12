package firewall

import (
	"encoding/json"
	"testing"
)

func TestNatPortForwardInterfaceUnmarshalMultipleSelected(t *testing.T) {
	data := []byte(`{
		"interface": {
			"wan": {"selected": 1, "value": "WAN"},
			"openvpn": {"selected": 1, "value": "OpenVPN"},
			"lan": {"selected": 0, "value": "LAN"}
		}
	}`)

	var rule NatPortForward
	if err := json.Unmarshal(data, &rule); err != nil {
		t.Fatalf("failed to unmarshal NAT port forward: %v", err)
	}

	if got, want := rule.Interface.String(), "openvpn,wan"; got != want {
		t.Fatalf("unexpected interface selection: got %q, want %q", got, want)
	}
}
