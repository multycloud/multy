package common

import (
	"regexp"
	"strings"
)

type FormatFunc func(s string) string

func LowercaseAlphanumericFormatFunc(s string) string {
	reg, err := regexp.Compile("[^a-z\\d]+")
	if err != nil {
		panic(err)
	}
	return reg.ReplaceAllString(strings.ToLower(s), "")
}

func AlphanumericFormatFunc(s string) string {
	reg, err := regexp.Compile("[^a-zA-Z\\d]+")
	if err != nil {
		panic(err)
	}
	return reg.ReplaceAllString(s, "")
}
