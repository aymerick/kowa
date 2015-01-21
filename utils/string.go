package utils

import (
	"math/rand"
	"strings"
	"time"
	"unicode"
)

var AlphaNumChars = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Generates a random alphanumeric string
func RandomAlphaNumString(size int) string {
	rand.Seed(time.Now().UTC().UnixNano())

	bytes := make([]rune, size)
	for i := range bytes {
		bytes[i] = AlphaNumChars[rand.Intn(len(AlphaNumChars))]
	}
	return string(bytes)
}

// Returns a string than can be used in an URL
func Urlify(str string) string {
	return strings.ToLower(UnicodeSanitize(strings.Replace(strings.TrimSpace(str), " ", "-", -1)))
}

// Borrowed from https://github.com/spf13/hugo/blob/master/helpers/path.go
func UnicodeSanitize(s string) string {
	source := []rune(s)
	target := make([]rune, 0, len(source))

	for _, r := range source {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '/' || r == '_' || r == '-' {
			target = append(target, r)
		}
	}

	return string(target)
}
