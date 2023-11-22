package strs

import "strings"

func StrIsEmail(str string) bool {
	s := strings.Split(str, "@")
	if len(s) < 2 {
		return false
	}
	if len(s[0]) < 1 {
		return false
	}
	if len(s[2]) < 1 {
		return false
	}

	return true
}
