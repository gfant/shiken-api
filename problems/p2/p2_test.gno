package p2

import (
	"testing"
)

func TestProblem(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{2, 0},
		{4, 2},
		{1000000, 78498},
	}

	for _, test := range tests {
		result := Problem(test.input)
		if result != test.expected {
			t.Errorf("Error in %q", test.input)
		}
	}
}
