package diagnostics

import (
	"encoding/json"
	"testing"
)

func TestCacheStatsUnmarshal(t *testing.T) {
	var empty CacheStats
	if err := json.Unmarshal([]byte(`[]`), &empty); err != nil {
		t.Fatalf("unmarshal empty array: %v", err)
	}
	if len(empty) != 0 {
		t.Fatalf("expected empty map, got %v", empty)
	}

	var populated CacheStats
	raw := `{"netflow_vtnet0":{"Pkts":23,"if":"vtnet0","SrcIPaddresses":2,"DstIPaddresses":2}}`
	if err := json.Unmarshal([]byte(raw), &populated); err != nil {
		t.Fatalf("unmarshal populated map: %v", err)
	}
	stat, ok := populated["netflow_vtnet0"]
	if !ok {
		t.Fatalf("expected key netflow_vtnet0 in %v", populated)
	}
	if stat.Pkts != 23 || stat.If != "vtnet0" || stat.SrcIPAddresses != 2 || stat.DstIPAddresses != 2 {
		t.Fatalf("unexpected CacheStat: %+v", stat)
	}
}
