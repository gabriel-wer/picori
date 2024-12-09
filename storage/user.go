package storage

import (
	"github.com/gabriel-wer/picori"
)

func (s *Sqlite) CheckUser(userFilter *picori.UserFilter) (*picori.User, error) {
	var user picori.User
	err := s.db.QueryRow("Select * from users where username = $1", userFilter.Username).Scan(&user)
	if err != nil {
		return &picori.User{}, err
	}

	return &user, nil
}

func (s *Sqlite) CreateUser(user picori.User) error {
	_, err := s.db.Exec("INSERT INTO users (id, username, created) VALUES ($1, $2, $3)", user.Id, user.Username, user.Created)
	if err != nil {
		return err
	}

	return nil
}
