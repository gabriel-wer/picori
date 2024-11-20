package main

import (
	"github.com/gabriel-wer/picori/server"
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

    server := server.NewServer(":6969", store)
    server.Start()
}
