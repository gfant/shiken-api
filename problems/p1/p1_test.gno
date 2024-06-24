package shikentest

import (
	"testing"
)

func TestProblem(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hola", "aloh"},
		{"mundo", "odnum"},
		{"", ""},
		{"a", "a"},
		{"ab", "ba"},
	}

	for _, test := range tests {
		result := Problem(test.input)
		if result != test.expected {
			t.Errorf("Error in %q", test.input)
		}
	}
}
