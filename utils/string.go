package utils

import "math/rand"

var AlphaNumChars = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Generates a random alphanumeric string
func RandomAlphaNumString(size int) string {
	bytes := make([]rune, size)
	for i := range bytes {
		bytes[i] = AlphaNumChars[rand.Intn(len(AlphaNumChars))]
	}
	return string(bytes)
}
