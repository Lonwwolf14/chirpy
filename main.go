package main

import (
	"database/sql"
	"log"
	"net/http"

	"example.com/chirpy/internal/app"
	"example.com/chirpy/internal/config"
	"example.com/chirpy/internal/database"
	"example.com/chirpy/internal/routes"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func initDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	const port = ":8080"
	router := mux.NewRouter()

	// Load Config
	apiConfig, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading configuration: %v", err)
	}

	// Initialize Database
	db, err := initDB(apiConfig.DbUrl)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	// Initialize App State
	queries := database.New(db)
	appState := &app.AppState{
		AppConfig: apiConfig,
		DB:        queries,
	}

	// Register Routes
	routes.Register(router, appState)

	// Start Server
	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}
	log.Printf("Starting server on %s", port)
	log.Fatal(srv.ListenAndServe())
}
