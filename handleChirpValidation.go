package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type validChirpRequest struct {
	Body string `json:"body"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type validChirpResponse struct {
	Body  string `json:"body"`
	Valid bool   `json:"valid"`
}

func cleanWord(word string) bool {
	word = strings.TrimSpace(word)
	word = strings.ToLower(word)
	if word == "kerfuffle" || word == "sharbert" || word == "fornax" {
		return false
	}
	return true
}

func (cfg *apiConfig) handleChirpValidation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(r.Body)
	requestBody := validChirpRequest{}
	err := decoder.Decode(&requestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{
			Error: "Invalid request body",
		})
		return
	}
	if len(requestBody.Body) > 140 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{
			Error: "Chirp is too long",
		})
		return
	}
	words := strings.Fields(requestBody.Body)
	for i, word := range words {
		if !cleanWord(word) {
			words[i] = "****"
		}
	}
	requestBody.Body = strings.Join(words, " ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(validChirpResponse{
		Body:  requestBody.Body,
		Valid: true,
	})

}
