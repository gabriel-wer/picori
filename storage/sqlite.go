package storage

import (
	"database/sql"
	"strings"

	"github.com/gabriel-wer/picori"
	_ "modernc.org/sqlite"
)

type Sqlite struct {
	db *sql.DB
}

func NewSqlite() *Sqlite {
	return &Sqlite{}
}

func (s *Sqlite) InitDB() {
	var err error
	s.db, err = sql.Open("sqlite", "database.db")
	if err != nil {
		panic(err)
	}
}

func (s *Sqlite) Close() error {
	err := s.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Sqlite) GetURL(url *picori.URL) error {
	err := s.db.QueryRow("SELECT longurl FROM url WHERE shorturl = $1", &url.ShortURL).Scan(&url.LongURL)
	if err != nil {
		return err
	}
	return nil
}

func (s *Sqlite) InsertURL(url picori.URL) error {
	_, err := s.db.Exec("INSERT INTO url (shorturl, longurl) VALUES ($1, $2)", url.ShortURL, url.LongURL)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			if err := url.Expand(); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

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
