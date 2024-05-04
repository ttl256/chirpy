package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ttl256/chirpy/internal/auth"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{Email: "", Password: ""}
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
	user := User{
		ID:    dbUser.ID,
		Email: dbUser.Email,
	}
	respondWithJSON(w, http.StatusOK, user)
}
