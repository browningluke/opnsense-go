package diagnostics

import "encoding/json"

// CacheStats unmarshals the netflow cache_stats response, which is a
// map keyed by interface name when populated but an empty JSON array
// ("[]") when netflow has collected nothing yet.
type CacheStats map[string]CacheStat

func (c *CacheStats) UnmarshalJSON(data []byte) error {
	var m map[string]CacheStat
	if err := json.Unmarshal(data, &m); err == nil {
		*c = m
		return nil
	}

	var empty []any
	if err := json.Unmarshal(data, &empty); err != nil {
		return err
	}
	*c = CacheStats{}
	return nil
}
