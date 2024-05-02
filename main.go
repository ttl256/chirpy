package main

import (
	"log"
	"net/http"
)

func main() {
	const tcpAddr = "0.0.0.0:8080"

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz", ReadinessHandler)

	s := http.Server{ //nolint: gosec // let me be
		Addr:    tcpAddr,
		Handler: mux,
	}

	log.Printf("Starting HTTP server on %s\n", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func ReadinessHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(http.StatusText(http.StatusOK)))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
