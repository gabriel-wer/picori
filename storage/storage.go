package storage

import (
	"github.com/gabriel-wer/picori/types"
)

type Storage interface {
    GetURL(*types.URL)
    InsertURL(types.URL) error
}
