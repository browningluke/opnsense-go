package api

import (
	"encoding/json"
	"testing"
)

func TestSelectedMapList_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected SelectedMapList
		wantErr  bool
	}{
		{
			name: "mixed bool and int selected values",
			input: `{
				"K1": {"selected": false, "value": "value1"},
				"K2": {"selected": true, "value": "value2"},
				"K3": {"selected": 1, "value": "value3"},
				"K4": {"selected": 0, "value": "value4"}
			}`,
			expected: SelectedMapList{"K2", "K3"},
			wantErr:  false,
		},
		{
			name:     "empty map",
			input:    `{}`,
			expected: SelectedMapList{},
			wantErr:  false,
		},
		{
			name:     "empty array fallback",
			input:    `[]`,
			expected: SelectedMapList{},
			wantErr:  false,
		},
		{
			name: "all selected false",
			input: `{
				"K1": {"selected": false, "value": "value1"},
				"K2": {"selected": 0, "value": "value2"}
			}`,
			expected: SelectedMapList{},
			wantErr:  false,
		},
		{
			name: "multiple selected items",
			input: `{
				"zebra": {"selected": true, "value": "z"},
				"alpha": {"selected": 1, "value": "a"},
				"beta": {"selected": true, "value": "b"}
			}`,
			expected: SelectedMapList{"alpha", "beta", "zebra"},
			wantErr:  false,
		},
		{
			name:    "invalid json",
			input:   `{invalid json}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s SelectedMapList
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

			if len(s) != len(tt.expected) {
				t.Errorf("Length mismatch: got %d, want %d", len(s), len(tt.expected))
				return
			}

			for i, item := range s {
				if item != tt.expected[i] {
					t.Errorf("Item %d: got %s, want %s", i, item, tt.expected[i])
				}
			}
		})
	}
}

func TestSelectedMapList_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    SelectedMapList
		expected string
	}{
		{
			name:     "multiple items",
			input:    SelectedMapList{"zebra", "alpha", "beta"},
			expected: "alpha,beta,zebra",
		},
		{
			name:     "single item",
			input:    SelectedMapList{"single"},
			expected: "single",
		},
		{
			name:     "empty list",
			input:    SelectedMapList{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.input.MarshalJSON()
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Unmarshal to get the actual string value from JSON
			var unmarshaled string
			err = json.Unmarshal(result, &unmarshaled)
			if err != nil {
				t.Errorf("Failed to unmarshal result: %v", err)
				return
			}

			if unmarshaled != tt.expected {
				t.Errorf("Got %s, want %s", unmarshaled, tt.expected)
			}
		})
	}
}

func TestSelectedMapList_String(t *testing.T) {
	tests := []struct {
		name     string
		input    SelectedMapList
		expected string
	}{
		{
			name:     "multiple items",
			input:    SelectedMapList{"alpha", "beta", "gamma"},
			expected: "alpha,beta,gamma",
		},
		{
			name:     "single item",
			input:    SelectedMapList{"single"},
			expected: "single",
		},
		{
			name:     "empty list",
			input:    SelectedMapList{},
			expected: "",
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

func TestSelectedMapListNL_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected SelectedMapListNL
		wantErr  bool
	}{
		{
			name: "mixed bool and int selected values",
			input: `{
				"K1": {"selected": false, "value": "value1"},
				"K2": {"selected": true, "value": "value2"},
				"K3": {"selected": 1, "value": "value3"}
			}`,
			expected: SelectedMapListNL{"K2", "K3"},
			wantErr:  false,
		},
		{
			name:     "empty array fallback",
			input:    `[]`,
			expected: SelectedMapListNL{},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s SelectedMapListNL
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

			if len(s) != len(tt.expected) {
				t.Errorf("Length mismatch: got %d, want %d", len(s), len(tt.expected))
				return
			}

			for i, item := range s {
				if item != tt.expected[i] {
					t.Errorf("Item %d: got %s, want %s", i, item, tt.expected[i])
				}
			}
		})
	}
}

func TestSelectedMapListNL_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    SelectedMapListNL
		expected string
	}{
		{
			name:     "multiple items with newlines",
			input:    SelectedMapListNL{"zebra", "alpha", "beta"},
			expected: "alpha\nbeta\nzebra",
		},
		{
			name:     "single item",
			input:    SelectedMapListNL{"single"},
			expected: "single",
		},
		{
			name:     "empty list",
			input:    SelectedMapListNL{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.input.MarshalJSON()
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Unmarshal to get the actual string value from JSON
			var unmarshaled string
			err = json.Unmarshal(result, &unmarshaled)
			if err != nil {
				t.Errorf("Failed to unmarshal result: %v", err)
				return
			}

			if unmarshaled != tt.expected {
				t.Errorf("Got %s, want %s", unmarshaled, tt.expected)
			}
		})
	}
}

func TestSelectedMapListNL_String(t *testing.T) {
	tests := []struct {
		name     string
		input    SelectedMapListNL
		expected string
	}{
		{
			name:     "multiple items with commas",
			input:    SelectedMapListNL{"alpha", "beta", "gamma"},
			expected: "alpha,beta,gamma",
		},
		{
			name:     "single item",
			input:    SelectedMapListNL{"single"},
			expected: "single",
		},
		{
			name:     "empty list",
			input:    SelectedMapListNL{},
			expected: "",
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
