package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/ttl256/chirpy/internal/auth"
	"github.com/ttl256/chirpy/internal/db"
)

func (cfg *apiConfig) membershipHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetAuth(r.Header, "ApiKey")
	if err != nil {
		log.Printf("error extracting token from header %#v: %s", r.Header, err)
		respondWithError(w, http.StatusUnauthorized, "error extracting token from header")
		return
	}
	if token != cfg.polkaAPIKey {
		log.Printf("api token is invalid. header %#v: %s", r.Header, err)
		respondWithError(w, http.StatusUnauthorized, "api token is invalid")
		return
	}
	type parametersData struct {
		UserID int `json:"user_id"`
	}
	type parameters struct {
		Event string         `json:"event"`
		Data  parametersData `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{
		Event: "",
		Data: parametersData{
			UserID: 0,
		},
	}
	if err = decoder.Decode(&params); err != nil {
		log.Printf("error decoding json: %s", err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("could not decode parameters: %s", err))
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusOK, nil)
		return
	}
	err = cfg.db.UpgradeMembership(params.Data.UserID)
	switch {
	case errors.Is(err, db.ErrNotExists):
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("could not find user with id %d", params.Data.UserID))
		return
	case err != nil:
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, nil)
}
