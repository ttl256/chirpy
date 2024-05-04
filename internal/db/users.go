package db

func (db *DB) CreateUser(email string) (User, error) {
	s, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	id := len(s.Users) + 1
	user := User{ID: id, Email: email}
	s.Users[id] = user
	err = db.writeDB(s)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
