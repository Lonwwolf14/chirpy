package app

import (
	"example.com/chirpy/internal/config"
	"example.com/chirpy/internal/database"
)

type AppState struct {
	AppConfig *config.ApiConfig
	DB        *database.Queries
}
