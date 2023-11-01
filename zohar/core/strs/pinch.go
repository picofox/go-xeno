package strs

import (
	"strings"
)

func ExtractPinchString(str string, left string, right string) (string, string, bool) {
	l := strings.Index(str, left)
	if l == -1 {
		return str, "", false
	}

	r := strings.Index(str[l+1:], right)
	if r == -1 {
		return str, "", false
	}

	r = r + l + 1

	if (l + 1) >= r {
		return str[0:l], "", true
	}

	return str[0:l], str[l+1 : r], true
}

func ExtractPinchStringList(str string, left string, right string, sep string, maxCount int) (string, []string, bool) {
	m, s, ok := ExtractPinchString(str, left, right)
	if !ok {
		return m, nil, ok
	}

	ret := strings.SplitN(s, sep, maxCount)
	return m, ret, ok
}
