package transform

import (
	"strconv"
	"strings"
)

func ToInt(value string) int {
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return -1
	}
	return int(v)
}

func TrimStrings(stringSlice []string) []string {
	var trimmed []string
	for _, str := range stringSlice {
		trimmed = append(trimmed, strings.TrimSpace(str))
	}
	return trimmed
}

func LowerStrings(stringSlice []string) []string {
	var lowered []string
	for _, str := range stringSlice {
		lowered = append(lowered, strings.ToLower(str))
	}
	return lowered
}

func NormalizeStrings(stringSlice []string) []string {
	return LowerStrings(TrimStrings(stringSlice))
}
