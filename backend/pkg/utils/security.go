package utils

import "strings"

func MaskSensitiveData(data string, prefixLen, suffixLen int, mask string) string {
	if prefixLen < 0 {
		prefixLen = 0
	}
	if suffixLen < 0 {
		suffixLen = 0
	}
	if mask == "" {
		mask = "***"
	}

	runes := []rune(data)
	length := len(runes)

	if length <= prefixLen+suffixLen+len([]rune(mask)) {
		return strings.Repeat("*", length)
	}

	if prefixLen > length {
		prefixLen = length 
	}
	if suffixLen > length-prefixLen {
		suffixLen = length - prefixLen
	}

	prefix := string(runes[:prefixLen])
	suffix := string(runes[length-suffixLen:])

	return prefix + mask + suffix
}