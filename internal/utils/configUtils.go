package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	GoogleOAuth2ClientId string `json:"googleOAuth2ClientId"`
	GoogleOAuth2Secret   string `json:"googleOAuth2Secret"`
	IdkBackendBaseUrl    string `json:"idkBackendBaseUrl"`
}

func LoadConfig() (*Config, error) {
	path := filepath.Join("configs", "appConfigs.json")
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	if err := json.NewDecoder(file).Decode(config); err != nil {
		return nil, err
	}
	return config, nil
}
