package db

import (
	"errors"
	"time"

	"github.com/ttl256/chirpy/internal/auth"
)

func (db *DB) CreateUser(email string, hash string) (User, error) {
	_, err := db.GetUserByEmail(email)
	switch {
	case err == nil:
		return User{}, ErrAlreadyExists
	case !errors.Is(err, ErrNotExists):
		return User{}, err
	}

	s, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	id := len(s.Users) + 1
	user := User{
		ID:        id,
		Email:     email,
		Password:  hash,
		ChirpyRed: false,
		RefreshToken: Token{
			Token:     "",
			ExpiresAt: time.Time{},
		},
	}
	s.Users[id] = user
	err = db.writeDB(s)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (db *DB) UpdateUser(id int, email, hash string) (User, error) {
	user, err := db.GetUserByID(id)
	if err != nil {
		return User{}, err
	}

	s, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user.Email = email
	user.Password = hash
	s.Users[id] = user

	err = db.writeDB(s)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (db *DB) GetUserByID(id int) (User, error) {
	s, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	user, ok := s.Users[id]
	if !ok {
		return User{}, ErrNotExists
	}
	return user, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	s, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	for _, user := range s.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return User{}, ErrNotExists
}

func (db *DB) CreateRefreshToken(id int) (string, error) {
	const n = 64
	user, err := db.GetUserByID(id)
	if err != nil {
		return "", err
	}
	token, err := auth.RandString64(n)
	if err != nil {
		return "", err
	}
	user.RefreshToken = Token{
		Token:     token,
		ExpiresAt: time.Now().UTC().Add(24 * 60 * time.Hour),
	}
	s, err := db.loadDB()
	if err != nil {
		return "", err
	}
	s.Users[id] = user
	err = db.writeDB(s)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (db *DB) UpgradeMembership(id int) error {
	user, err := db.GetUserByID(id)
	if err != nil {
		return err
	}
	user.ChirpyRed = true
	s, err := db.loadDB()
	if err != nil {
		return err
	}
	s.Users[user.ID] = user
	err = db.writeDB(s)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetUserByRefreshToken(token string) (User, error) {
	s, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	for _, user := range s.Users {
		if user.RefreshToken.Token == token {
			return user, nil
		}
	}
	return User{}, ErrNotExists
}

func (db *DB) RevokeToken(token string) error {
	s, err := db.loadDB()
	if err != nil {
		return err
	}
	for _, user := range s.Users {
		if user.RefreshToken.Token == token {
			user.RefreshToken = Token{
				Token:     "",
				ExpiresAt: time.Time{},
			}
			s.Users[user.ID] = user
			err = db.writeDB(s)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return ErrNotExists
}
