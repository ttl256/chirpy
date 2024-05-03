package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	if code >= http.StatusInternalServerError {
		log.Printf("responding with 5xx error: %s", msg)
	}
	respondWithJSON(w, code, errorResponse{Error: msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error encoding JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	_, err = w.Write(data)
	if err != nil {
		log.Printf("error writing JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
