package db

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  sync.Mutex
}

func New(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  sync.Mutex{},
	}
	err := db.ensureDB()
	return db, err
}

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

func (db *DB) createDB() error {
	s := Schema{
		Chirps: map[int]Chirp{},
	}
	return db.writeDB(s)
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

func (db *DB) loadDB() (Schema, error) {
	schema := Schema{
		Chirps: map[int]Chirp{},
	}
	db.mux.Lock()
	defer db.mux.Unlock()
	data, err := os.ReadFile(db.path)
	if err != nil {
		return Schema{}, err
	}
	err = json.Unmarshal(data, &schema)
	if err != nil {
		return Schema{}, err
	}
	return schema, nil
}

func (db *DB) writeDB(schema Schema) error {
	data, err := json.Marshal(schema)
	if err != nil {
		return err
	}
	db.mux.Lock()
	defer db.mux.Unlock()
	if err = os.WriteFile(db.path, data, 0600); err != nil {
		return err
	}
	return nil
}
