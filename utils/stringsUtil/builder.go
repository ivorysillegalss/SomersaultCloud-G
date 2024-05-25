package stringsUtil

import "strings"

func Concat(str ...string) string {
	var builder strings.Builder
	for i := range str {
		builder.WriteString(str[i])
	}
	return builder.String()
}
