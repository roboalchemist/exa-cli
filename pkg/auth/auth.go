package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type AuthConfig struct {
	APIKey string `json:"api_key"`
}

// GetAPIKey returns the API key from env var or config file.
// Priority: EXA_API_KEY env var → ~/.exa-auth.json
func GetAPIKey() (string, error) {
	if key := os.Getenv("EXA_API_KEY"); key != "" {
		return key, nil
	}

	config, err := loadAuth()
	if err != nil {
		return "", fmt.Errorf("EXA_API_KEY not set and no config file found.\nRun 'exa auth' to configure or set EXA_API_KEY environment variable")
	}
	if config.APIKey != "" {
		return config.APIKey, nil
	}
	return "", fmt.Errorf("no valid authentication found")
}

// GetBaseURL returns the API base URL, with env var override.
func GetBaseURL() string {
	if url := os.Getenv("EXA_API_URL"); url != "" {
		return url
	}
	return "https://api.exa.ai"
}

func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".exa-auth.json"), nil
}

func loadAuth() (*AuthConfig, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config AuthConfig
	return &config, json.Unmarshal(data, &config)
}

func SaveAuth(config AuthConfig) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}
