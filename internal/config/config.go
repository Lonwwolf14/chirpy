package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync/atomic"
)

type ApiConfig struct {
	fileserverHits atomic.Int64
	DbUrl          string `json:"db_url"`
}

const configFileName = ".gatorconfig.json" // Update this path if needed

// Read reads the configuration from the user's home directory
func Read() (ApiConfig, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ApiConfig{}, err
	}

	// Use filepath.Join to handle platform-specific path separators
	configFilePath := filepath.Join(homeDir, configFileName)

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return ApiConfig{}, err
	}

	var configData ApiConfig
	err = json.Unmarshal(data, &configData)
	if err != nil {
		return ApiConfig{}, err
	}

	return configData, nil
}

// GetFileServerHits returns the current number of hits, atomically
func (cfg *ApiConfig) GetFileServerHits() int64 {
	return cfg.fileserverHits.Load()
}

// AddHit increments the file server hit counter atomically
func (cfg *ApiConfig) AddHit() {
	cfg.fileserverHits.Add(1)
}

// ResetFileServerHits resets the file server hit counter atomically
func (cfg *ApiConfig) ResetFileServerHits() {
	cfg.fileserverHits.Store(0)
}
