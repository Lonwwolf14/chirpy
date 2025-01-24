package handlers

import (
	"fmt"
	"net/http"

	"example.com/chirpy/internal/app"
)

func MiddlewareMetricsInc(s *app.AppState, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Incrementing file server hit count...")
		s.AppConfig.AddHit()
		next.ServeHTTP(w, r)
	})
}
