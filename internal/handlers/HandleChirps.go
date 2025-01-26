package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"example.com/chirpy/internal/app"
	"example.com/chirpy/internal/auth"
	"example.com/chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Chirp struct {
	Body   string    `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}

func HandleChirp(s *app.AppState, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		//Checking the JWT Token
		tokenString, err := auth.GetBearerToken(r.Header)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		//Validate the token to get the user ID
		userID, err := auth.ValidateJWT(tokenString, s.AppConfig.TokenSecret)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		//Decoding the JSON body
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		chirp := &Chirp{}
		err = decoder.Decode(chirp)
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
			UserID:    userID,
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
			chirpID, createdAt.Format(time.RFC3339), updatedAt.Format(time.RFC3339), chirp.Body, userID,
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

func HandleChirpById(s *app.AppState, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	chirpId := vars["chirp_id"]
	if chirpId == "" {
		http.Error(w, "Missing chirp_id parameter", http.StatusBadRequest)
		return
	}
	if chirpId == "" {
		http.Error(w, "Missing chirp_id parameter", http.StatusBadRequest)
		return
	}
	id, err := uuid.Parse(chirpId)
	if err != nil {
		http.Error(w, "Invalid chirp_id parameter", http.StatusBadRequest)
		return
	}
	chirp, err := s.DB.GetChirp(r.Context(), id)
	if err != nil {
		http.Error(w, "Error getting chirp", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chirp)

}

func DeleteChirpByID(s *app.AppState, w http.ResponseWriter, r *http.Request) {
	// Extract the Bearer token from the request header
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		fmt.Printf("Failed to extract Bearer token: %v\n", err)
		http.Error(w, "Failed to extract authentication token. Please ensure you are logged in.", http.StatusUnauthorized)
		return
	}

	// Validate the extracted JWT token
	userID, err := auth.ValidateJWT(token, s.AppConfig.TokenSecret)
	if err != nil {
		fmt.Printf("Failed to validate JWT token: %v\n", err)
		http.Error(w, "Invalid or expired token. Please log in again.", http.StatusUnauthorized)
		return
	}

	// Parse the chirp ID from the request URL and validate its presence
	chirpID, err := uuid.Parse(mux.Vars(r)["chirp_id"])
	if err != nil {
		fmt.Printf("Invalid chirp ID format: %v\n", err)
		http.Error(w, "Invalid chirp ID format. Please provide a valid chirp ID.", http.StatusBadRequest)
		return
	}

	//Get userID using chirpID
	chirp, err := s.DB.GetChirp(r.Context(), chirpID)
	if err != nil {
		fmt.Printf("Failed to get chirp from database: %v\n", err)
		http.Error(w, "Failed to get chirp. It may not exist.", http.StatusNotFound)
		return
	}
	chirpUserID := chirp.UserID

	if userID != chirpUserID {
		fmt.Printf("User ID mismatch: %v != %v\n", userID, chirpUserID)
		http.Error(w, "You do not have permission to delete this chirp.", http.StatusForbidden)
		return
	}

	// Attempt to delete the chirp from the database
	err = s.DB.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID:     chirpID,
		UserID: userID,
	})
	if err != nil {
		fmt.Printf("Failed to delete chirp from database: %v\n", err)
		http.Error(w, "Failed to delete chirp. It may not exist, or you may not have permission to delete it.", http.StatusForbidden)
		return
	}

	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := fmt.Sprintf(`{"message": "Chirp deleted successfully.", "user_id": "%s", "chirp_id": "%s"}`, userID.String(), chirpID.String())
	w.Write([]byte(response))
}
