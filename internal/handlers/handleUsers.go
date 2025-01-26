package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"example.com/chirpy/internal/app"
	"example.com/chirpy/internal/auth"
	"example.com/chirpy/internal/database"
	"github.com/google/uuid"
)

type user struct {
	Email string `json:"email"`
	Pass  string `json:"password"`
}

type WebhookData struct {
	UserID string `json:"user_id"`
}

type WebhookRequest struct {
	Event string      `json:"event"`
	Data  WebhookData `json:"data"`
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

	//User Creation
	userID := uuid.New()
	currentTime := time.Now()
	createdAt := sql.NullTime{Time: currentTime, Valid: true}
	updatedAt := sql.NullTime{Time: currentTime, Valid: true}
	passwd, err := auth.HashPassword(requestBody.Pass)
	if err != nil {
		fmt.Printf("Error hashing password: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error hashing password"))
		return
	}
	_, err = s.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        userID,            // Pass the generated UUID
		Email:     requestBody.Email, // Email from the request
		CreatedAt: createdAt,         // Current timestamp
		UpdatedAt: updatedAt,         // Current timestamp
		Password:  passwd,            // Password from the request
	})
	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error creating user(might already exists)"))
		return

	}

	//Token Generation
	token, err := auth.MakeJWT(userID, s.AppConfig.TokenSecret, 24*time.Hour)
	if err != nil {
		fmt.Printf("Error generating token: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error generating token"))
		return
	}

	//Response Struct
	response := struct {
		ID        uuid.UUID    `json:"id"`
		Email     string       `json:"email"`
		CreatedAt sql.NullTime `json:"created_at"`
		UpdatedAt sql.NullTime `json:"updated_at"`
		Token     string       `json:"token"`
	}{
		ID:        userID,
		Email:     requestBody.Email,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Token:     token,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("user created"))
	json.NewEncoder(w).Encode(response)

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

func UpdateUser(s *app.AppState, w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	decoder := json.NewDecoder(r.Body)
	requestBody := user{}
	err = decoder.Decode(&requestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userID, err := auth.ValidateJWT(token, s.AppConfig.TokenSecret)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	newPassword, err := auth.HashPassword(requestBody.Pass)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	_, err = s.DB.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:        userID,
		Email:     requestBody.Email,
		Password:  newPassword,
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := struct {
		ID        uuid.UUID    `json:"id"`
		Email     string       `json:"email"`
		CreatedAt sql.NullTime `json:"created_at"`
		UpdatedAt sql.NullTime `json:"updated_at"`
	}{
		ID:        userID,
		Email:     requestBody.Email,
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User updated"))
	json.NewEncoder(w).Encode(response)

}

func HandlePolkaWebhook(s *app.AppState, w http.ResponseWriter, r *http.Request) {
	//authenticating user
	key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if key != os.Getenv("POLKA_KEY") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//decoder the body
	decoder := json.NewDecoder(r.Body)
	requestBody := WebhookRequest{}
	err = decoder.Decode(&requestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//Check the event
	if requestBody.Event != "user.upgraded" {
		http.Error(w, "Invalid event", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(requestBody.Data.UserID)
	if err != nil {
		return
	}
	err = s.DB.UpgradeToRed(r.Context(), userID)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User upgraded to Red"))

}
