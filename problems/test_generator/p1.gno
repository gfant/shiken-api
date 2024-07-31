package shikentest

import (
	"testing"
)

func TestProblem(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{ 
		{"oi4XgJ", "JgX4io"},
		{"vbM", "Mbv"},
		{"", ""},
		{"HX2v", "v2XH"},
		{"EbRY", "YRbE"},
		{"r9", "9r"},
		{"e6", "6e"},
		{"S5", "5S"},
		{"Ovfc", "cfvO"},
		{"yU4MXpDZ3", "3ZDpXM4Uy"},
	}

	for _, test := range tests {
		result := Problem(test.input)
		if result != test.expected {
			t.Errorf("Error in %q", test.input)
		}
	}
}
