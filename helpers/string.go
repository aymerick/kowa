package helpers

import (
	"math/rand"
	"net/url"
	"strings"
	"time"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
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

func Pathify(s string) string {
	return strings.ToLower(NormalizeToPath(strings.Replace(strings.TrimSpace(s), " ", "-", -1)))
}

// Returns a string than can be used in an URL
func Urlify(str string) string {
	// escape unicode letters
	parsedUri, err := url.Parse(str)
	if err != nil {
		panic(err)
	}

	return parsedUri.String()
}

// Normalize unicode string to a string that can be a file or URL path
func NormalizeToPath(str string) string {
	isNotOk := func(r rune) bool {
		isOk := (r == 35) || // '#'
			(r == 45) || // '-'
			(r == 46) || // '.'
			(r == 47) || // '/'
			((r >= 48) && (r <= 57)) || // '0'..'9'
			((r >= 65) && (r <= 90)) || // 'A'..'Z'
			((r >= 97) && (r <= 122)) || // 'a'..'z'
			(r == 95) // '_'

		return !isOk
	}

	t := transform.Chain(norm.NFKD, transform.RemoveFunc(isNotOk))

	result, _, _ := transform.String(t, str)

	return result
}

// Normalise unicode string to a string that can be a path part
func NormalizeToPathPart(str string) string {
	isNotOk := func(r rune) bool {
		isOk := ((r >= 48) && (r <= 57)) || // '0'..'9'
			((r >= 65) && (r <= 90)) || // 'A'..'Z'
			((r >= 97) && (r <= 122)) // 'a'..'z'

		return !isOk
	}

	t := transform.Chain(norm.NFKD, transform.RemoveFunc(isNotOk))

	result, _, _ := transform.String(t, str)

	return result
}

// Returns true is string has one of given prefixes
func HasOnePrefix(s string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}

	return false
}
