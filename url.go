package picori

import (
	"crypto/sha1"
	"errors"
	"fmt"
)

const HASH_SIZE = 9

type URL struct {
	ShortURL string `json:"shorturl"`
	LongURL  string `json:"longurl"`
}

func (url *URL) Expand() error {
	if len(url.ShortURL) == 0 {
		err := errors.New("You already have an expanded URL.")
		return err
	}
	return nil
}

func (url *URL) Shorten() error {
	if len(url.ShortURL) == 0 || len(url.LongURL) == 0 {
		err := errors.New("You already have a shortened URL.")
		return err
	}

	hash := sha1.Sum([]byte(url.LongURL))
	url.ShortURL = fmt.Sprintf("%x", hash)[:6]
	return nil
}
