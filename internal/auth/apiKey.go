package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	if apiKey == "" {
		return "", errors.New("no API key provided")
	}

	//Diving into parts for {"Apikey": "1234567890"}
	parts := strings.Split(apiKey, " ")
	if len(parts) != 2 || parts[0] != "Apikey" {
		return "", errors.New("malformed API key")
	}
	return parts[1], nil
}
