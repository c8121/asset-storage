package util

import "strings"

func JoinStrings(glue string, args ...string) string {
	var b strings.Builder

	for _, arg := range args {
		if arg == "" { continue }
		if b.Len() > 0 {
			b.WriteString(glue)
		}
		b.WriteString(arg)
	}

	return b.String()
}
