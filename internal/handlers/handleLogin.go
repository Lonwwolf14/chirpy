package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"example.com/chirpy/internal/app"
	"example.com/chirpy/internal/auth"
	"github.com/google/uuid"
)

type LoginRequest struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds *int   `json:"expires_in_seconds"`
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
		log.Printf("Database error: %v", err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	log.Printf("User from DB - ID: %v", user.ID)
	err = auth.CheckPasswordHash(loginRequest.Password, user.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Credentials are invalid"))
		return
	}

	//Handle the expiration time
	const defaultExpiration = 3600
	expirationTime := defaultExpiration
	if loginRequest.ExpiresInSeconds != nil {
		if *loginRequest.ExpiresInSeconds < defaultExpiration && *loginRequest.ExpiresInSeconds > 0 {
			expirationTime = *loginRequest.ExpiresInSeconds
		}
	}
	//generate the JWT token
	token, err := auth.MakeJWT(user.ID, s.AppConfig.TokenSecret, time.Duration(expirationTime)*time.Second)
	if err != nil {
		log.Printf("Token generation error: %v", err)
		http.Error(w, "Error generating JWT", http.StatusInternalServerError)
		return
	}
	log.Printf("Generated token for ID: %v", user.ID)
	//Create the response struct
	response := struct {
		ID        uuid.UUID `json:"id"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Token     string    `json:"token"`
	}{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		Token:     token,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
