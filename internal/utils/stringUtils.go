package utils

import (
	"crypto/rand"
	"regexp"
)

func GenerateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, n) // slice to hold the random bytes

	// Generate random numbers and map them to the indices of the letters slice
	for i := range bytes {
		if _, err := rand.Read(bytes[i : i+1]); err != nil {
			return ""
		}
		bytes[i] = letters[bytes[i]%byte(len(letters))]
	}

	return string(bytes)
}

func RemoveWhiteSpaceFromString(s string) string {
	// Compile regex to match all whitespace
	re := regexp.MustCompile(`\s+`)
	// Replace all whitespace with nothing
	noWhitespaceString := re.ReplaceAllString(s, "")
	return noWhitespaceString
}
