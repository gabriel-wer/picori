package storage

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/gabriel-wer/picori"
	_ "github.com/lib/pq"
)

type Postgre struct {
	db *sql.DB
}

func NewPostgre() *Postgre {
	return &Postgre{}
}

func (p *Postgre) GetURL(url *picori.URL) error {
	err := p.db.QueryRow("SELECT longurl FROM url WHERE shorturl = $1", &url.ShortURL).Scan(&url.LongURL)
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgre) InsertURL(url picori.URL) error {
	_, err := p.db.Exec("INSERT INTO url (shorturl, longurl) VALUES ($1, $2)", url.ShortURL, url.LongURL)
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

func (p *Postgre) InitDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=picori sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PORT"))

	var err error
	p.db, err = sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	//TODO: Change to Schema file
	_, err = p.db.Exec("CREATE TABLE IF NOT EXISTS url (shorturl TEXT, longurl TEXT, PRIMARY KEY(shorturl))")
	if err != nil {
		panic(err)
	}
}

func (p *Postgre) Close() error {
	err := p.db.Close()
	if err != nil {
		return err
	}

	return nil
}
