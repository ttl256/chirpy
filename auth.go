package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/ttl256/chirpy/internal/auth"
)

func (cfg *apiConfig) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearer(r.Header)
		if err != nil {
			log.Printf("error extracting JWT from header %#v: %s", r.Header, err)
			respondWithError(w, http.StatusUnauthorized, "error extracting JWT from header")
			return
		}
		subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
		if err != nil {
			log.Printf("error validating JWT %s: %s", token, err)
			respondWithError(w, http.StatusUnauthorized, "error validating JWT")
			return
		}
		id, err := strconv.Atoi(subject)
		if err != nil {
			log.Printf("error parsing user ID: %s", err)
			respondWithError(w, http.StatusInternalServerError, "error parsing user ID")
			return
		}
		ctx := context.WithValue(r.Context(), subjectID{}, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
