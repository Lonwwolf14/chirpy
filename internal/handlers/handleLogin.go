package handlers

import (
	"encoding/json"
	"net/http"

	"example.com/chirpy/internal/app"
	"example.com/chirpy/internal/auth"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func HandleLogin(s *app.AppState, w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var loginRequest LoginRequest
	err := decoder.Decode(&loginRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := s.DB.GetUserPassword(r.Context(), loginRequest.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	err = auth.CheckPasswordHash(loginRequest.Password, user.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Credentials are invalid"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))

}
