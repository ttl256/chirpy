package main

import (
	"log"
	"net/http"
)

func main() {
	const tcpAddr = "0.0.0.0:8080"

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))

	s := http.Server{ //nolint: gosec // let me be
		Addr:    tcpAddr,
		Handler: mux,
	}

	log.Printf("Starting HTTP server on %s\n", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
