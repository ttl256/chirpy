package main

import (
	"log"
	"net/http"

	"github.com/ttl256/chirpy/internal/db"
)

func main() {
	const tcpAddr = "0.0.0.0:8080"

	db, err := db.New("database.json")
	if err != nil {
		log.Fatal(err)
	}
	cfg := &apiConfig{
		fileserverHits: 0,
		db:             db,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /api/reset", cfg.metricsResetHandler)
	mux.HandleFunc("GET /api/chirps", cfg.getChirpsHandler)
	mux.HandleFunc("POST /api/chirps", cfg.postChirpHandler)

	mux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)

	s := http.Server{ //nolint: gosec // let me be
		Addr:    tcpAddr,
		Handler: mux,
	}

	log.Printf("Starting HTTP server on %s\n", s.Addr)
	if err = s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

type apiConfig struct {
	fileserverHits int
	db             *db.DB
}
