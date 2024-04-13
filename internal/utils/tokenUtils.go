package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// TokenData wraps both access and refresh tokens
type TokenData struct {
	JwtToken string `json:"jwtToken"`
}

// SaveToken saves the token data to a file
func SaveToken(jwtToken string) error {
	data := TokenData{
		JwtToken: jwtToken,
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	secretsFilePath := GetAbsoluteHomeDirectoryPath([]string{".idk", "credentials"})

	dirPath := filepath.Dir(secretsFilePath)
	if err := os.MkdirAll(dirPath, 0777); err != nil {
		return err
	}

	return os.WriteFile(secretsFilePath, bytes, 0600) // Use 0600 to restrict access to the file owner
}

// LoadToken loads the token data from a file
func LoadToken() (string, error) {
	secretsFilePath := GetAbsoluteHomeDirectoryPath([]string{".idk", "credentials"})

	bytes, err := os.ReadFile(secretsFilePath)
	if err != nil {
		return "", err
	}
	var data TokenData
	if err := json.Unmarshal(bytes, &data); err != nil {
		return "", err
	}
	return data.JwtToken, nil
}

func ClearToken() error {
	secretsFilePath := GetAbsoluteHomeDirectoryPath([]string{".idk", "credentials"})

	return os.Remove(secretsFilePath)
}

func RemoveCredentialsFromCommand(command string) string {
	// Define patterns to match credentials.
	patterns := []string{
		`(\-\-|\-)?(password|passwd|pwd|pass|api[_\-]?key|token|secret|signature)(=|\s+)[^ ]+`, // Matches most CLI tool credential patterns
		`(\S+)=['"]?[^'"\s]+['"]?`, // Matches key=value pairs, optionally enclosed in quotes
		`Bearer\s+[^ ]+`,           // Matches 'Bearer ' followed by a token, commonly used in HTTP authorization
		`Basic\s+[^ ]+`,            // Matches 'Basic ' followed by a base64 string, used in HTTP authorization
	}

	// Replace each pattern found with an empty string or a placeholder.
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		command = re.ReplaceAllString(command, "")
	}

	// Clean up any extra spaces left by removed credentials
	command = strings.Join(strings.Fields(command), " ")

	return command
}
