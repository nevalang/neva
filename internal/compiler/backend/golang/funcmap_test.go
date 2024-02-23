package golang

import (
	"testing"
)

func Test_handleSpecialChars(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "const_ref_sender",
			input:    "$greeting",
			expected: "_greeting",
		},
		{
			name:     "normal port addr",
			input:    "foo:bar",
			expected: "foo_bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handleSpecialChars(tt.input)
			if result != tt.expected {
				t.Errorf("handleSpecialChars() = %v, want %v", result, tt.expected)
			}
		})
	}
}
