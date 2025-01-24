package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"example.com/chirpy/internal/app"
	"example.com/chirpy/internal/database"
	"github.com/google/uuid"
)

type user struct {
	Email string `json:"email"`
}

func HandleUserCreation(s *app.AppState, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(r.Body)
	requestBody := user{}
	err := decoder.Decode(&requestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//User
	userID := uuid.New()
	currentTime := time.Now()
	createdAt := sql.NullTime{Time: currentTime, Valid: true}
	updatedAt := sql.NullTime{Time: currentTime, Valid: true}

	_, err = s.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        userID,            // Pass the generated UUID
		Email:     requestBody.Email, // Email from the request
		CreatedAt: createdAt,         // Current timestamp
		UpdatedAt: updatedAt,         // Current timestamp
	})
	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error creating user(might already exists)"))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("user created"))
	if createdAt.Valid && updatedAt.Valid {
		w.Write([]byte(fmt.Sprintf("User Id: %s\nEmail: %s\nCreated At: %s\nUpdated At: %s", userID, requestBody.Email, createdAt.Time.Format(time.RFC3339), updatedAt.Time.Format(time.RFC3339))))
	} else {
		w.Write([]byte("Error: CreatedAt or UpdatedAt is not valid"))
	}

}

func HandleUsersDeletion(s *app.AppState, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	_, err := s.DB.DeleteUsers(r.Context())
	if err != nil {
		fmt.Printf("Error deleting users: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error deleting users"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All users deleted"))
}
