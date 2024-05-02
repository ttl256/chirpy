package main

import (
	"log"
	"net/http"
)

func handlerReadiness(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(http.StatusText(http.StatusOK)))
	// TODO how to handle this?
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	const addr = "127.0.0.1"
	const port = "8080"

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz", handlerReadiness)

	corsMux := middlewareCors(mux)

	srv := http.Server{
		Addr:              addr + ":" + port,
		Handler:           corsMux,
		ReadHeaderTimeout: 0,
	}

	log.Printf("Started server on %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
