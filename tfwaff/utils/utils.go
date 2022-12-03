package utils

import (
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func IsPascalCase(input string) bool {
	return regexp.MustCompile(`^[A-Z][a-z]+(?:[A-Z][a-z]+)*$`).MatchString(input)
}

func GetSnakeCase(input string) string {
	if !IsPascalCase(input) {
		return input
	}
	return strings.ToLower(strings.TrimPrefix(regexp.MustCompile(`([A-Z][a-z]+)`).ReplaceAllString(input, "_$1"), "_"))
}

func GetKebabCase(input string) string {
	if !IsPascalCase(input) {
		return input
	}
	return strings.ToLower(strings.TrimPrefix(regexp.MustCompile(`([A-Z][a-z]+)`).ReplaceAllString(input, "-$1"), "-"))
}

func GetTitleCase(input string) string {
	lower := strings.ToLower(regexp.MustCompile(`_|-`).ReplaceAllString(input, " "))
	if !regexp.MustCompile(`^[a-z]+(\s?[a-z]+)*$`).MatchString(lower) {
		return input
	}
	return cases.Title(language.Und, cases.NoLower).String(lower)
}
