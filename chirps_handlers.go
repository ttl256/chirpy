package main

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Chirp struct {
	ID       int    `json:"id"`
	AuthorID int    `json:"author_id"`
	Body     string `json:"body"`
}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, _ *http.Request) {
	dbChirps, err := cfg.db.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	chirps := make([]Chirp, 0, len(dbChirps))
	for _, chirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:       chirp.ID,
			AuthorID: chirp.AuthorID,
			Body:     chirp.Body,
		})
	}

	slices.SortStableFunc(chirps, func(a, b Chirp) int {
		return cmp.Compare(a.ID, b.ID)
	})

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) chirpByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("chirp_id"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("invalid ID: %s", err))
		return
	}
	dbChirp, err := cfg.db.GetChirpByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	chirp := Chirp{
		ID:       dbChirp.ID,
		AuthorID: dbChirp.AuthorID,
		Body:     dbChirp.Body,
	}
	respondWithJSON(w, http.StatusOK, chirp)
}

func (cfg *apiConfig) deleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("chirp_id"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("invalid ID: %s", err))
		return
	}
	dbChirp, err := cfg.db.GetChirpByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	authorID, ok := r.Context().Value(subjectID{}).(int)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "could not infer type int on field subject")
		return
	}

	if authorID != dbChirp.AuthorID {
		respondWithError(w, http.StatusForbidden, "cannot delete someone else's chirp")
		return
	}

	err = cfg.db.DeleteChirp(dbChirp.ID)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	authorID, ok := r.Context().Value(subjectID{}).(int)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "could not infer type int on field subject")
		return
	}

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{Body: ""}
	if err := decoder.Decode(&params); err != nil {
		log.Printf("error decoding json: %s", err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("could not decode parameters: %s", err))
		return
	}

	body, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	chirp, err := cfg.db.CreateChirp(authorID, body)
	if err != nil {
		log.Printf("error creating chirp: %s", err)
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:       chirp.ID,
		AuthorID: authorID,
		Body:     chirp.Body,
	})
}

func validateChirp(body string) (string, error) {
	const maxLen = 140
	if utf8.RuneCountInString(body) > maxLen {
		return "", errors.New("Chirp is too long")
	}

	profanes := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	body = filterWords(body, profanes, "****")
	return body, nil
}

func filterWords(s string, w map[string]struct{}, sub string) string {
	words := strings.Fields(s)
	for i, word := range words {
		if _, ok := w[strings.ToLower(word)]; ok {
			words[i] = sub
		}
	}
	return strings.Join(words, " ")
}
