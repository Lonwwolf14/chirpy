package handlers

import (
	"net/http"

	"example.com/chirpy/internal/app"
)

func MiddlewareMetricsInc(s *app.AppState, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.AppConfig.AddHit()
		next.ServeHTTP(w, r)
	})
}
