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

func LowercaseAlphanumericAndDashFormatFunc(s string) string {
	s = strings.ToLower(s)
	// remove all illegal chars
	s = regexp.MustCompile("[^-a-z\\d]+").ReplaceAllString(s, "")
	// remove dashes and numbers in the beginning of the string
	s = regexp.MustCompile("^[-0-9]+").ReplaceAllString(s, "")
	// remove dashes from the end of the string
	s = regexp.MustCompile("-+$").ReplaceAllString(s, "")
	return s
}
