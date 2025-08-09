package api

import (
	"encoding/json"
	"strconv"
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
	}

	SelectedMap allows the JSON library to unmarshal that map into a string containing only
	the key that is selected (i.e. "K2", in the example above).
*/

type SelectedMap string

func (s *SelectedMap) UnmarshalJSON(data []byte) error {
	var dataMap map[string]struct {
		Value    string `json:"value"`
		Selected any    `json:"selected"`
	}

	err := json.Unmarshal(data, &dataMap)
	if err != nil {
		// some responses return a list of selected items instead of a map
		var listValues []struct {
			Value    string `json:"value"`
			Selected any    `json:"selected"`
		}
		listErr := json.Unmarshal(data, &listValues)
		if listErr == nil {
			for idx, v := range listValues {
				if selectedBool, ok := v.Selected.(bool); ok {
					if selectedBool {
						*s = SelectedMap(strconv.Itoa(idx))
					}
				}
				if selectedInt, ok := v.Selected.(float64); ok {
					if selectedInt == 1 {
						*s = SelectedMap(strconv.Itoa(idx))
					}
				}
			}
			return nil
		}
		return err
	}

	// Find selected element
	for k, v := range dataMap {
		// If bool
		if selectedBool, ok := v.Selected.(bool); ok {
			if selectedBool {
				*s = SelectedMap(k)
			}
		}

		// If float64
		if selectedInt, ok := v.Selected.(float64); ok {
			if selectedInt == 1 {
				*s = SelectedMap(k)
			}
		}
	}

	return nil
}

func (s *SelectedMap) String() string {
	return string(*s)
}
