package main

import (
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		log.Printf("Hits: %d", cfg.fileserverHits)
		next.ServeHTTP(w, r)
	})
}

const templ = `<html>

<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>

</html>`

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(fmt.Sprintf(templ, cfg.fileserverHits)))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (cfg *apiConfig) metricsResetHandler(w http.ResponseWriter, _ *http.Request) {
	cfg.fileserverHits = 0
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
