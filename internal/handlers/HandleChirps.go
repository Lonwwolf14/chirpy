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

func HandleChirp(s *app.AppState, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
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

	case http.MethodGet:
		chirps, err := s.DB.GetAllChirps(r.Context())
		if err != nil {
			http.Error(w, "Error getting chirps", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(chirps)
		return

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return

	}

}
