package utils

import "strings"

func ToSnakeCase(str string) string {
	var result strings.Builder
	runes := []rune(str)

	for i, r := range runes {
		if i > 0 && 'A' <= r && r <= 'Z' {
			// Check if previous character was lowercase or if next character is lowercase
			prevIsLower := i > 0 && 'a' <= runes[i-1] && runes[i-1] <= 'z'
			nextIsLower := i < len(runes)-1 && 'a' <= runes[i+1] && runes[i+1] <= 'z'

			if prevIsLower || nextIsLower {
				result.WriteByte('_')
			}
		}

		if 'A' <= r && r <= 'Z' {
			result.WriteRune(r - 'A' + 'a')
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}
