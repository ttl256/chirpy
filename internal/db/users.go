package db

import "errors"

func (db *DB) CreateUser(email string, hash string) (User, error) {
	if _, err := db.GetUserByEmail(email); !errors.Is(err, ErrNotExists) {
		return User{}, err
	}

	s, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	id := len(s.Users) + 1
	user := User{ID: id, Email: email, Password: hash}
	s.Users[id] = user
	err = db.writeDB(s)
	if err != nil {
		return User{}, err
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
