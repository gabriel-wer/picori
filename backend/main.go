package main

import (
	"github.com/gabriel-wer/picori/api"
	"github.com/gabriel-wer/picori/storage"
)

func main() {
    store := storage.NewPostgre()
    store.InitDB()
    defer func() {
        if err := store.Close(); err != nil {
            panic(err)
        }
    }()

    server := api.NewServer(":3000", store)
    server.Start()
}
