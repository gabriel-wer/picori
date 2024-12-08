package storage

import (
	"github.com/gabriel-wer/picori"
)

type Storage interface {
	GetURL(*picori.URL) error
	InsertURL(picori.URL) error
}
