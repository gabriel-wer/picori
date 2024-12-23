package storage

import (
	"github.com/gabriel-wer/picori"
	"strings"
)

func (s *Sqlite) GetURL(url *picori.URL) error {
	err := s.db.QueryRow("SELECT longurl FROM url WHERE shorturl = $1", &url.ShortURL).Scan(&url.LongURL)
	if err != nil {
		return err
	}
	return nil
}

func (s *Sqlite) ListURL() ([]picori.URL, error) {
	var urls []picori.URL
	rows, err := s.db.Query("SELECT longurl, shorturl FROM url")
	if err != nil {
		return []picori.URL{}, err
	}
    defer rows.Close()

    for rows.Next() {
        var url picori.URL
        err = rows.Scan(&url.LongURL, &url.ShortURL)
        if err != nil {
            return []picori.URL{}, err
        }
        urls = append(urls, url)
    }

	return urls, nil
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
