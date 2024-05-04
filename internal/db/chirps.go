package db

import "fmt"

func (db *DB) CreateChirp(body string) (Chirp, error) {
	s, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	id := len(s.Chirps) + 1
	c := Chirp{ID: id, Body: body}
	s.Chirps[id] = c
	err = db.writeDB(s)
	if err != nil {
		return Chirp{}, err
	}
	return c, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	s, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}
	chirps := make([]Chirp, 0, len(s.Chirps))
	for _, i := range s.Chirps {
		chirps = append(chirps, i)
	}
	return chirps, nil
}

func (db *DB) GetChirpByID(id int) (Chirp, error) {
	s, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	chirp, ok := s.Chirps[id]
	if !ok {
		return Chirp{}, fmt.Errorf("no chirp with ID %d", id)
	}
	return chirp, nil
}
