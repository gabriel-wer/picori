package storage

import (
	"database/sql"
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
