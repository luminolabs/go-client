package utils

import "strings"

func ContainsStringFromArray(source string, subStringArray []string) bool {
	for i := 0; i < len(subStringArray); i++ {
		if strings.Contains(source, subStringArray[i]) {
			return true
		}
	}
	return false
}
