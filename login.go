package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ttl256/chirpy/internal/auth"
)

const hour = 60 * 60

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		Expiration int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{Email: "", Password: "", Expiration: 0}
	if err := decoder.Decode(&params); err != nil {
		log.Printf("error decoding json: %s", err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("could not decode parameters: %s", err))
		return
	}

	dbUser, err := cfg.db.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "no such user")
		return
	}

	err = auth.CheckPasswordHash(params.Password, dbUser.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "password is invalid")
		return
	}

	if params.Expiration <= 0 || params.Expiration > hour {
		params.Expiration = hour
	}

	token, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, time.Second*time.Duration(params.Expiration))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating JWT")
		return
	}

	tokenRefresh, err := cfg.db.CreateRefreshToken(dbUser.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating refresh token")
		return
	}

	user := struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}{
		User: User{
			ID:    dbUser.ID,
			Email: dbUser.Email,
		},
		Token:        token,
		RefreshToken: tokenRefresh,
	}
	respondWithJSON(w, http.StatusOK, user)
}

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	tokenRefresh, err := auth.GetBearer(r.Header)
	if err != nil {
		log.Printf("error extracting refresh token from header %#v: %s", r.Header, err)
		respondWithError(w, http.StatusUnauthorized, "error extracting refresh token from header")
		return
	}

	dbUser, err := cfg.db.GetUserByRefreshToken(tokenRefresh)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if dbUser.RefreshToken.ExpiresAt.Before(time.Now().UTC()) {
		respondWithError(w, http.StatusUnauthorized, "refresh token is expired")
		return
	}

	token, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating JWT")
		return
	}

	response := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	tokenRefresh, err := auth.GetBearer(r.Header)
	if err != nil {
		log.Printf("error extracting refresh token from header %#v: %s", r.Header, err)
		respondWithError(w, http.StatusUnauthorized, "error extracting refresh token from header")
		return
	}

	err = cfg.db.RevokeToken(tokenRefresh)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJSON(w, http.StatusOK, nil)
}
