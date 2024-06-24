package shikentest

import "strings"

func Problem(s string) string {
	result := ""
	for _, c := range s {
		if c != 'a' && c != 'e' && c != 'i' && c != 'o' && c != 'u' {
			result = result + string(c)
		}
	}
	return result
}
