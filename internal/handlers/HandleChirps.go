package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"example.com/chirpy/internal/app"
	"example.com/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	Body   string    `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}

func HandleCreateChirp(s *app.AppState, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	chirp := &Chirp{}
	err := decoder.Decode(chirp)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	chirpID := uuid.New()
	currentTime := time.Now()
	createdAt := currentTime
	updatedAt := currentTime
	_, err = s.DB.CreateChirp(r.Context(), database.CreateChirpParams{
		ID:        chirpID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserId,
	})
	if err != nil {
		http.Error(w, "Error creating chirp", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Write the response using w.Write([]byte(...))
	w.Write([]byte(fmt.Sprintf(
		"Chirp created successfully\nId: %s\nCreatedAt: %s\nUpdatedAt: %s\nBody: %s\nUserId: %s\n",
		chirpID, createdAt.Format(time.RFC3339), updatedAt.Format(time.RFC3339), chirp.Body, chirp.UserId,
	)))
}
