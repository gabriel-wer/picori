package main

import (
	"github.com/gabriel-wer/picori/api"
	"github.com/gabriel-wer/picori/storage"
)

func main() {
	store := storage.NewSqlite()
	store.InitDB()
	defer func() {
		if err := store.Close(); err != nil {
			panic(err)
		}
	}()

	server := api.NewServer(store)
	server.Start()
}
