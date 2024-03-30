package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
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
