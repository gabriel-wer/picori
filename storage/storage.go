package storage

import (
	"github.com/gabriel-wer/picori/types"
)

type Storage interface {
    GetURL(*types.URL) error
    InsertURL(types.URL) error
}
