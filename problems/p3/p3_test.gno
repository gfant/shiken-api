package p1

import (
	"testing"
)

func TestProblem(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"gnoland", "gnlnd"},
		{"hola", "hl"},
		{"mundo", "mnd"},
		{"jorin", "jrn"},
		{"this is amazing", "ths s mzng"},
	}

	for _, test := range tests {
		result := Problem(test.input)
		if result != test.expected {
			t.Errorf("Error in %q", test.input)
		}
	}
}
