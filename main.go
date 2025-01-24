package main

import (
	"database/sql"
	"log"
	"net/http"

	"example.com/chirpy/internal/app"
	"example.com/chirpy/internal/config"
	"example.com/chirpy/internal/database"
	"example.com/chirpy/internal/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	const port = ":8080"
	router := mux.NewRouter()

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
	registerRoutes(router, appState) // Changed mux to router
	srv := &http.Server{
		Addr:    port,
		Handler: router, // Changed mux to router
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

func registerRoutes(router *mux.Router, appState *app.AppState) { // Changed mux to router
	router.Handle("/app/", handlers.MiddlewareMetricsInc(appState, http.StripPrefix("/app", http.FileServer(http.Dir("."))))).Methods("GET")
	router.HandleFunc("/api/healthz", handleReadiness).Methods("GET")
	router.HandleFunc("/admin/metrics", wrapHandler(appState, handlers.HandleMetrics)).Methods("GET")
	router.HandleFunc("/api/validate_chirp", wrapHandler(appState, handlers.HandleChirpValidation)).Methods("POST")
	router.HandleFunc("/api/users", wrapHandler(appState, handlers.HandleUserCreation)).Methods("POST")
	router.HandleFunc("/admin/reset", wrapHandler(appState, handlers.HandleUsersDeletion)).Methods("POST")
	router.HandleFunc("/api/chirps", wrapHandler(appState, handlers.HandleChirp)).Methods("POST")
	router.HandleFunc("/api/chirps", wrapHandler(appState, handlers.HandleChirp)).Methods("GET")
	router.HandleFunc("/api/chirps/{chirp_id}", wrapHandler(appState, handlers.HandleChirpById)).Methods("GET")
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
