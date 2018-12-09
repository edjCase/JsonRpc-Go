package common

import (
	"fmt"
	"strings"
)

//IsNullOrWhitespace will return true if the string is null, empty strings or just whitespace,
//otherwise it will return false
func IsNullOrWhitespace(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

//SplitAndClean will split the string and remove all entries that are
//null, empty strings or just whitespace
func SplitAndClean(s string, seperator string) []string {
	parts := strings.Split(s, seperator)

	var cleanedParts []string
	for _, part := range parts {
		if !IsNullOrWhitespace(part) {
			cleanedParts = append(cleanedParts, part)
		}
	}
	return cleanedParts
}

func NestedError(err error, message string) error {
	if err == nil {
		panic("Nested error cannot be nil")
	}
	return fmt.Errorf("%s, Inner error: %s", message, err.Error())
}
