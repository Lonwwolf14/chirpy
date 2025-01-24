package main

import (
	"log"
	"net/http"

	"example.com/chirpy/handlers"
)

func main() {
	mux := http.NewServeMux()
	const port = ":8080"
	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}
	apiConfig := handlers.ApiConfig{}

	mux.Handle("/app/", apiConfig.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("/api/healthz", readiness)
	mux.HandleFunc("/admin/metrics", apiConfig.HandleMetrics)
	mux.HandleFunc("/admin/reset", apiConfig.HandleReset)
	mux.HandleFunc("/api/validate_chirp", apiConfig.HandleChirpValidation)
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
