package api

import (
	"encoding/json"
	"sort"
	"strings"
)

/*
	OPNsense responses to some queries with json data that looks like:
	"some_key" : {
		"K1": {
			"selected": 0 (or false),
			"value": "...",
		},
		"K2": {
			"selected": 1 (or true),
			"value": "...",
		},
		"K3": {
			"selected": true (or 1),
			"value": "...",
		},
		"K4": {
			"selected": false (or 0),
			"value": "...",
		},
	}

	SelectedMapList allows the JSON library to unmarshal that map into a string containing only
	the key(s) that is/are selected (i.e. "K2,K3", in the example above).
*/

type SelectedMapList []string

func (s *SelectedMapList) UnmarshalJSON(data []byte) error {
	var dataMap map[string]struct {
		Value    string `json:"value"`
		Selected any    `json:"selected"`
	}

	err := json.Unmarshal(data, &dataMap)
	if err != nil {
		return err
	}

	*s = SelectedMapList{}
	// Find selected element
	for k, v := range dataMap {
		// If bool
		if selectedBool, ok := v.Selected.(bool); ok {
			if selectedBool {
				*s = append(*s, k)
			}
		}

		// If float64
		if selectedInt, ok := v.Selected.(float64); ok {
			if selectedInt == 1 {
				*s = append(*s, k)
			}
		}
	}
	sort.Strings(*s)

	return nil
}

func (s *SelectedMapList) MarshalJSON() ([]byte, error) {
	// Ensure list is sorted
	sort.Strings(*s)
	str := strings.Join(*s, ",")
	return json.Marshal(str)
}

func (s *SelectedMapList) String() string {
	return strings.Join(*s, ",")
}
