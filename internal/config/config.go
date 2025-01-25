package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/joho/godotenv"
)

type ApiConfig struct {
	fileserverHits atomic.Int64
	DbUrl          string `json:"db_url"`
	TokenSecret    string
}

const configFileName = ".gatorconfig.json"

func Read() (*ApiConfig, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configFilePath := filepath.Join(homeDir, configFileName)

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	var configData ApiConfig
	err = json.Unmarshal(data, &configData)
	if err != nil {
		return nil, err
	}
	configData.TokenSecret = os.Getenv("JWT_SECRET")
	return &configData, nil
}

func (cfg *ApiConfig) GetFileServerHits() int64 {
	return cfg.fileserverHits.Load()
}

func (cfg *ApiConfig) AddHit() {
	cfg.fileserverHits.Add(1)
}

func (cfg *ApiConfig) ResetFileServerHits() {
	cfg.fileserverHits.Store(0)
}
