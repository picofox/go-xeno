package strs

import (
	"crypto/md5"
	"strings"
)

var sStr64 string = "%abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789#"

func CreateSampleStringWithMD5(totalLen int, begin string, end string) (string, []byte) {
	str := CreateSampleString(totalLen, begin, end)
	ba := md5.Sum([]byte(str))
	return str, ba[:]
}

func CreateSampleString(totalLen int, begin string, end string) string {
	if totalLen <= 0 {
		return ""
	}
	var ss strings.Builder
	strLen := totalLen - len(begin) - len(end)
	remainLen := strLen
	if len(begin) > 0 {
		ss.WriteString(begin)
	}
	for remainLen > 0 {
		if remainLen > len(sStr64) {
			ss.WriteString(sStr64)
			remainLen -= len(sStr64)
		} else {
			ss.WriteString(sStr64[0:remainLen])
			remainLen -= remainLen
		}
	}
	if len(end) > 0 {
		ss.WriteString(end)
	}
	return ss.String()
}

func CreateSampleText(s string, loop int, maxLen int, extra string) string {
	if maxLen < len(s) {
		return ""
	}
	var ss strings.Builder
	for i := 0; i < loop; i++ {
		if maxLen > 0 {
			if ss.Len()+len(s) > maxLen {
				remain := maxLen - ss.Len()
				ss.WriteString(s[0:remain])
			} else {
				ss.WriteString(s)
			}
		} else {
			ss.WriteString(s)
		}
	}

	if maxLen > 0 {
		if ss.Len() != maxLen {
			panic("Create Sample String Failed")
		}
	}

	if len(extra) > 0 {
		ss.WriteString(extra)
	}

	return ss.String()
}
