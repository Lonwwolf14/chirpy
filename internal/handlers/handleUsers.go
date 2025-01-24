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

func HandleUsers(s *app.AppState, w http.ResponseWriter, r *http.Request) {
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
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("user created"))

}
