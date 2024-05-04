package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/ttl256/chirpy/internal/db"
)

func main() {
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	const dbName = "database.json"
	if *debug {
		err := os.Remove(dbName)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				log.Fatal(err)
			}
			log.Printf("database %q does not exist, nothing to remove, continue", dbName)
		} else {
			log.Printf("database %q is removed, continue", dbName)
		}
	}

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
	mux.HandleFunc("GET /api/chirps/{chirp_id}", cfg.chirpByIDHandler)
	mux.HandleFunc("POST /api/chirps", cfg.createChirpHandler)

	mux.HandleFunc("POST /api/users", cfg.createUserHandler)

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
