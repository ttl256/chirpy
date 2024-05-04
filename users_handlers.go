package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ttl256/chirpy/internal/auth"
)

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
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

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error hashing password: %s", err))
		return
	}

	user, err := cfg.db.CreateUser(params.Email, hash)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, User{
		ID:    user.ID,
		Email: user.Email,
	},
	)
}
