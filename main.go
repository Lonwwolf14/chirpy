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
	mux := http.NewServeMux()
	const port = ":8080"
	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	apiConfig, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", apiConfig.DbUrl)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close() // Ensure the DB connection is closed when the program ends
	queries := database.New(db)

	appState := &app.AppState{
		AppConfig: &apiConfig,
		DB:        queries,
	}

	mux.Handle("/app/", handlers.MiddlewareMetricsInc(appState, http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("/api/healthz", readiness)
	mux.HandleFunc("/admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleMetrics(appState, w, r)
	})
	mux.HandleFunc("/admin/reset", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleReset(appState, w, r)
	})
	mux.HandleFunc("/api/validate_chirp", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleChirpValidation(appState, w, r)
	})
	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleUsers(appState, w, r)
	})
	log.Fatal(srv.ListenAndServe())

}

func readiness(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
