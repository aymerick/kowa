package utils

import (
	"math/rand"
	"time"
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
