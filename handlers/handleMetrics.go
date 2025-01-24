package handlers

import (
	"fmt"
	"net/http"
)

func (cfg *ApiConfig) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	htmlContent := fmt.Sprintf(`
		<html>
  		<body>
    		<h1>Welcome, Chirpy Admin</h1>
    		<p>Chirpy has been visited %d times!</p>
  		</body>
		</html>`, cfg.fileserverHits.Load())

	w.Write([]byte(htmlContent))

}
