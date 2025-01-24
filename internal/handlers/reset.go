package handlers

import (
	"net/http"

	"example.com/chirpy/internal/app"
)

func HandleReset(s *app.AppState, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	s.AppConfig.GetFileServerHits()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
