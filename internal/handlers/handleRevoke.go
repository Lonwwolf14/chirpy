package handlers

import (
	"fmt"
	"net/http"

	"example.com/chirpy/internal/app"
	"example.com/chirpy/internal/auth"
)

func HandleRevoke(s *app.AppState, w http.ResponseWriter, r *http.Request) {

	refresh_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		fmt.Print(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	_, err = s.DB.RevokeRefreshToken(r.Context(), refresh_token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Refresh token revoked"))

}
