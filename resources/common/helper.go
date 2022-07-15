package common

import (
	"github.com/multycloud/multy/validate"
	"hash/fnv"
	"log"
	"math/rand"
	"regexp"
)

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func RemoveSpecialChars(a string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(a, "")
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

// UniqueId generates a stable string composed of prefix+suffix and a 4 char hash.
// Prefix can be any size but will be sliced if bigger than 16 chars. Suffix can have 4 chars at most.
// Returns a string with at most 24 chars.
func UniqueId(prefix string, suffix string, formatFunc FormatFunc) string {
	if len(suffix) > 10 {
		validate.LogInternalError("suffix must be shorter than 10 chars")
	}
	result := ""
	formattedPrefix := formatFunc(prefix)
	maxPrefixLen := 20 - len(suffix)
	if len(formattedPrefix) > maxPrefixLen {
		result += formattedPrefix[:maxPrefixLen] + GenerateHash(prefix)
	} else {
		result += formattedPrefix
	}

	result += suffix
	result += GenerateHash(prefix + suffix)

	return result
}

func GenerateHash(s string) string {
	result := ""
	h := fnv.New32a()
	_, err := h.Write([]byte(s))
	if err != nil {
		panic(err)
	}

	hashSum := h.Sum32()
	var letters = []rune("abcdefghijklmnopqrstuvwxyz123456789")

	for i := 0; i < 4; i++ {
		idx := int(hashSum >> (i * 8) & 0xFF)
		result += string(letters[idx%len(letters)])
	}
	return result
}
