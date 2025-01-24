package main

import (
	"database/sql"
	"log"
	"net/http"

	"example.com/chirpy/internal/app"
	"example.com/chirpy/internal/config"
	"example.com/chirpy/internal/database"
	"example.com/chirpy/internal/handlers"
	_ "github.com/lib/pq"
)

func main() {
	const port = ":8080"
	mux := http.NewServeMux()

	apiConfig, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading configuration: %v", err)
	}

	db, err := initDB(apiConfig.DbUrl)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	queries := database.New(db)
	appState := &app.AppState{
		AppConfig: apiConfig,
		DB:        queries,
	}
	registerRoutes(mux, appState)
	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}
	log.Printf("Starting server on %s", port)
	log.Fatal(srv.ListenAndServe())
}

func initDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func registerRoutes(mux *http.ServeMux, appState *app.AppState) {
	mux.Handle("/app/", handlers.MiddlewareMetricsInc(appState, http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("/api/healthz", handleReadiness)
	mux.HandleFunc("/admin/metrics", wrapHandler(appState, handlers.HandleMetrics))
	mux.HandleFunc("/api/validate_chirp", wrapHandler(appState, handlers.HandleChirpValidation))
	mux.HandleFunc("/api/users", wrapHandler(appState, handlers.HandleUserCreation))
	mux.HandleFunc("/admin/reset", wrapHandler(appState, handlers.HandleUsersDeletion))
}

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	respondWithPlainText(w, http.StatusOK, http.StatusText(http.StatusOK))
}

func wrapHandler(appState *app.AppState, handlerFunc func(*app.AppState, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlerFunc(appState, w, r)
	}
}

func respondWithPlainText(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	w.Write([]byte(message))
}
