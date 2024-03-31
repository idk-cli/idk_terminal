package configs

import (
	"embed"
	"encoding/json"
	"io/fs"
)

//go:embed appConfigs.json
var configFS embed.FS // Embedding the specific file

type Config struct {
	IdkBackendBaseUrl string `json:"idkBackendBaseUrl"`
}

func LoadConfig() (*Config, error) {
	data, err := fs.ReadFile(configFS, "appConfigs.json")
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return config, nil
}
