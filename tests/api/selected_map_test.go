package api

import (
	"encoding/json"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestSelectedMap_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected api.SelectedMap
		wantErr  bool
	}{
		{
			name: "map format with bool selected",
			input: `{
				"K1": {"selected": false, "value": "value1"},
				"K2": {"selected": true, "value": "value2"},
				"K3": {"selected": false, "value": "value3"}
			}`,
			expected: api.SelectedMap("K2"),
			wantErr:  false,
		},
		{
			name: "map format with int selected",
			input: `{
				"K1": {"selected": 0, "value": "value1"},
				"K2": {"selected": 1, "value": "value2"},
				"K3": {"selected": 0, "value": "value3"}
			}`,
			expected: api.SelectedMap("K2"),
			wantErr:  false,
		},
		{
			name: "list format with bool selected",
			input: `[
				{"selected": false, "value": "value1"},
				{"selected": true, "value": "value2"},
				{"selected": false, "value": "value3"}
			]`,
			expected: api.SelectedMap("value2"),
			wantErr:  false,
		},
		{
			name: "multiple selected items in map (one will be selected)",
			input: `{
				"K1": {"selected": true, "value": "value1"},
				"K2": {"selected": true, "value": "value2"}
			}`,
			expected: api.SelectedMap(""), // We'll check that it's non-empty instead
			wantErr:  false,
		},
		{
			name: "no selected items in map",
			input: `{
				"K1": {"selected": false, "value": "value1"},
				"K2": {"selected": 0, "value": "value2"}
			}`,
			expected: api.SelectedMap(""),
			wantErr:  false,
		},
		{
			name: "no selected items in list",
			input: `[
				{"selected": false, "value": "value1"},
				{"selected": false, "value": "value2"}
			]`,
			expected: api.SelectedMap(""),
			wantErr:  false,
		},
		{
			name:     "empty map",
			input:    `{}`,
			expected: api.SelectedMap(""),
			wantErr:  false,
		},
		{
			name:     "empty list",
			input:    `[]`,
			expected: api.SelectedMap(""),
			wantErr:  false,
		},
		{
			name:    "invalid json",
			input:   `{invalid json}`,
			wantErr: true,
		},
		{
			name: "mixed selected types in map",
			input: `{
				"K1": {"selected": false, "value": "value1"},
				"K2": {"selected": 1, "value": "value2"},
				"K3": {"selected": true, "value": "value3"}
			}`,
			expected: api.SelectedMap(""), // We'll check that it's non-empty instead
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s api.SelectedMap
			err := json.Unmarshal([]byte(tt.input), &s)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Special handling for cases where map iteration order is non-deterministic
			if tt.name == "multiple selected items in map (one will be selected)" {
				if s == "" {
					t.Errorf("Expected a selected item but got empty string")
				}
				// Should be either K1 or K2
				if s != "K1" && s != "K2" {
					t.Errorf("Got %s, expected either K1 or K2", s)
				}
				return
			}

			if tt.name == "mixed selected types in map" {
				if s == "" {
					t.Errorf("Expected a selected item but got empty string")
				}
				// Should be either K2 or K3
				if s != "K2" && s != "K3" {
					t.Errorf("Got %s, expected either K2 or K3", s)
				}
				return
			}

			if s != tt.expected {
				t.Errorf("Got %s, want %s", s, tt.expected)
			}
		})
	}
}

func TestSelectedMap_String(t *testing.T) {
	tests := []struct {
		name     string
		input    api.SelectedMap
		expected string
	}{
		{
			name:     "non-empty string",
			input:    api.SelectedMap("test-value"),
			expected: "test-value",
		},
		{
			name:     "empty string",
			input:    api.SelectedMap(""),
			expected: "",
		},
		{
			name:     "string with special characters",
			input:    api.SelectedMap("test-value_123"),
			expected: "test-value_123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.String()
			if result != tt.expected {
				t.Errorf("Got %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestSelectedMap_RoundTrip(t *testing.T) {
	// Test that we can create a api.SelectedMap and convert it to string properly
	original := "test-key"
	selected := api.SelectedMap(original)

	result := selected.String()
	if result != original {
		t.Errorf("Round trip failed: got %s, want %s", result, original)
	}
}
