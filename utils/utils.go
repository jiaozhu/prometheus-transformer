package utils

import (
	"strings"
	"unicode"
)

func CamelCaseToSnakeCase(input string) string {
	lastUnderscore := strings.LastIndex(input, "_")
	if lastUnderscore == -1 {
		return input
	}

	prefix := input[:lastUnderscore+1]
	camelCasePortion := input[lastUnderscore+1:]

	var result string
	var lastPos int
	for pos, char := range camelCasePortion {
		if unicode.IsUpper(char) && pos > 0 {
			result += camelCasePortion[lastPos:pos] + "_"
			lastPos = pos
		}
	}
	result += camelCasePortion[lastPos:]

	return prefix + result
}
