package routes

import (
	"net/http"

	"example.com/chirpy/internal/app"
	"example.com/chirpy/internal/handlers"
	"github.com/gorilla/mux"
)

func Register(router *mux.Router, appState *app.AppState) {
	router.Handle("/app/", handlers.MiddlewareMetricsInc(appState, http.StripPrefix("/app", http.FileServer(http.Dir("."))))).Methods("GET")
	router.HandleFunc("/api/healthz", handleReadiness).Methods("GET")
	router.HandleFunc("/admin/metrics", wrapHandler(appState, handlers.HandleMetrics)).Methods("GET")
	router.HandleFunc("/api/validate_chirp", wrapHandler(appState, handlers.HandleChirpValidation)).Methods("POST")
	router.HandleFunc("/api/users", wrapHandler(appState, handlers.HandleUserCreation)).Methods("POST")
	router.HandleFunc("/admin/reset", wrapHandler(appState, handlers.HandleUsersDeletion)).Methods("DELETE")
	router.HandleFunc("/api/chirps", wrapHandler(appState, handlers.HandleChirp)).Methods("POST", "GET")
	router.HandleFunc("/api/login", wrapHandler(appState, handlers.HandleLogin)).Methods("POST")
	router.HandleFunc("/api/refresh", wrapHandler(appState, handlers.RefreshHandler)).Methods("POST")
	router.HandleFunc("/api/revoke", wrapHandler(appState, handlers.HandleRevoke)).Methods("POST")
	router.HandleFunc("/api/users", wrapHandler(appState, handlers.UpdateUser)).Methods("PUT")
	router.HandleFunc("/api/chirps/{chirp_id}", wrapHandler(appState, handlers.HandleChirpById)).Methods("GET")
	router.HandleFunc("/api/chirps/{chirp_id}", wrapHandler(appState, handlers.DeleteChirpByID)).Methods("DELETE")
}

func wrapHandler(appState *app.AppState, handlerFunc func(*app.AppState, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlerFunc(appState, w, r)
	}
}

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	respondWithPlainText(w, http.StatusOK, http.StatusText(http.StatusOK))
}

func respondWithPlainText(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	w.Write([]byte(message))
}
