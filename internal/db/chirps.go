package db

import "fmt"

func (db *DB) CreateChirp(authorID int, body string) (Chirp, error) {
	s, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	id := len(s.Chirps) + 1
	c := Chirp{
		ID:       id,
		AuthorID: authorID,
		Body:     body,
	}
	s.Chirps[id] = c
	err = db.writeDB(s)
	if err != nil {
		return Chirp{}, err
	}
	return c, nil
}

func (db *DB) DeleteChirp(id int) error {
	s, err := db.loadDB()
	if err != nil {
		return err
	}
	delete(s.Chirps, id)
	err = db.writeDB(s)
	if err != nil {
		return err
	}
	return nil
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

func (db *DB) GetChirpsByAuthor(id int) ([]Chirp, error) {
	s, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}
	chirps := make([]Chirp, 0, len(s.Chirps))
	for _, i := range s.Chirps {
		if i.AuthorID == id {
			chirps = append(chirps, i)
		}
	}
	return chirps, nil
}
