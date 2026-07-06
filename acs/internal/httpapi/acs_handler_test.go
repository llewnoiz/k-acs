package httpapi

import "testing"

func TestParseUpTime(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"", 0},
		{"123", 123},
		{"0", 0},
		{"abc", 0},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := parseUpTime(test.input)
			if result != test.expected {
				t.Errorf("parseUpTime(%q) = %d; want %d", test.input, result, test.expected)
			}
		})
	}
}
