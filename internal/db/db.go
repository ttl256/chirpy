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

func (db *DB) createDB() error {
	s := Schema{
		Chirps: map[int]Chirp{},
		Users:  map[int]User{},
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
		Users:  map[int]User{},
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
	data, err := json.MarshalIndent(schema, "", "  ")
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
