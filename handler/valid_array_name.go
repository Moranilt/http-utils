package handler

import "strings"

func isValidNameArray(name string) (int, bool) {
	lastIndex := len(name) - 1
	if lastIndex == -1 || name[lastIndex] != ']' {
		return -1, false
	}
	open := strings.Index(name, "[")
	return open, open != -1 && name[lastIndex] == ']' && lastIndex-open == 1
}

func extractArrayName(name string) (string, bool) {
	i, valid := isValidNameArray(name)
	if !valid {
		return name, false
	}

	return name[:i], true
}

func isValidSubName(name string) (int, bool) {
	lastIndex := len(name) - 1
	if lastIndex == -1 || name[lastIndex] != ']' {
		return -1, false
	}
	open := strings.Index(name, "[")
	return open, open != -1 && name[lastIndex] == ']' && lastIndex-open > 1
}

func extractSubName(name string) (string, string, bool) {
	i, valid := isValidSubName(name)
	if !valid {
		return name, "", false
	}
	return name[:i], name[i+1 : len(name)-1], true
}
