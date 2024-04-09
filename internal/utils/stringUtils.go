package utils

import (
	"crypto/rand"
	"regexp"
	"sort"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
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

func FilterByPrefix(stringsSlice []string, prefix string) []string {
	var filtered []string
	for _, str := range stringsSlice {
		if strings.HasPrefix(str, prefix) {
			filtered = append(filtered, str)
		}
	}
	return filtered
}

func FindMostRelevantStringFromArr(arr []string, s string) string {
	matches := fuzzy.RankFindNormalizedFold(s, arr)

	// Sort matches by their rank. Higher rank means a closer match.
	sort.Sort(sort.Reverse(matches))

	// Extract just the matched strings in sorted order.
	var sortedMatches []string
	for _, match := range matches {
		sortedMatches = append(sortedMatches, match.Target)
	}

	if len(matches) == 0 {
		return ""
	}

	return sortedMatches[0]
}
