package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/gabriel-wer/picori/storage"
)

type Server struct {
	ListenAddr string `env:"LISTENADDR" envDefault:":6969"`
	store      *storage.Sqlite
}

func NewServer(store *storage.Sqlite) Server {
	s := Server{
		store: store,
	}

	err := env.Parse(&s)
	if err != nil {
		panic(err)
	}

	return s
}

func (s *Server) Start() {
	r := NewRouter(CORS, Logging)

	r.Post("/v1/shorten", s.handleShorten)
	r.Post("/v1/expand", s.handleExpand)
	r.Post("/v1/login", s.handleLogin)
	r.Get("/v1/{url}", s.handleRedirect)
	r.Get("/v1/list", Authentication(s.handleList, s.store))

	log.Fatal(http.ListenAndServe(s.ListenAddr, r))
}

func (s *Server) serveStatic() {
	directoryPath := "./frontend/"

	_, err := os.Stat(directoryPath)
	if os.IsNotExist(err) {
		fmt.Printf("Directory '%s' not found.\n", directoryPath)
		return
	}

	fileServer := http.FileServer(http.Dir(directoryPath))

	http.Handle("/", fileServer)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
