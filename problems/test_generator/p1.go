package main

import (
	"math/rand"
	"unicode/utf8"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var TemplateP1Tests = `package shikentest

import (
	"testing"
)

func TestProblem(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{ {{range .}}
		{{"{"}}"{{.Test}}", "{{.Result}}"{{"}"}},{{end}}
	}

	for _, test := range tests {
		result := Problem(test.input)
		if result != test.expected {
			t.Errorf("Error in %q", test.input)
		}
	}
}
`

func SolutionP1(s string) string {
	size := len(s)
	buf := make([]byte, size)
	for start := 0; start < size; {
		r, n := utf8.DecodeRuneInString(s[start:])
		start += n
		utf8.EncodeRune(buf[size-start:], r)
	}
	return string(buf)
}

func TestGeneratorP1() (string, string) {
	size := rand.Intn(10)
	testVal := ""
	resultTest := ""
	for i := 0; i < size; i++ {
		testVal += string(charset[rand.Intn(len(charset))])
		resultTest = SolutionP1(testVal)
	}
	return testVal, resultTest
}
