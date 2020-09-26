package tfguard

import "strings"

// stringInSlice checks if string is contained in a slice
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// stringStartsInSlice checks if any items in the list match the start of the string
func stringStartsInSlice(a string, list []string) bool {
	for _, b := range list {
		if strings.HasPrefix(a, b) {
			return true
		}
	}
	return false
}
