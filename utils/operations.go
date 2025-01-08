package utils

import "strings"

// ContainsStringFromArray checks if any string from the array is a substring of the source.
// Used for error message matching and validation operations.
func ContainsStringFromArray(source string, subStringArray []string) bool {
	for i := 0; i < len(subStringArray); i++ {
		if strings.Contains(source, subStringArray[i]) {
			return true
		}
	}
	return false
}
