package handlers

import (
	"fmt"
	"net/http"

	"example.com/chirpy/internal/app"
)

func HandleMetrics(s *app.AppState, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	htmlContent := fmt.Sprintf(`
		<html>
  		<body>
    		<h1>Welcome, Chirpy Admin</h1>
    		<p>Chirpy has been visited %d times!</p>
  		</body>
		</html>`, s.AppConfig.GetFileServerHits())

	w.Write([]byte(htmlContent))

}
