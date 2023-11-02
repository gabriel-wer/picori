package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gabriel-wer/picori/storage"
	"github.com/gabriel-wer/picori/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)


type Server struct {
    listenAddr string
    store storage.Storage
}

func NewServer(listenAddr string, store storage.Storage) *Server {
    return &Server{
        listenAddr: listenAddr,
        store: store,
    }
}

func (s *Server) Start() error{

    corsMiddleware := cors.New(cors.Options{
        AllowedOrigins: []string{"*"},
        AllowedMethods: []string{"GET", "POST"},
        AllowedHeaders: []string{"*"},
    })
    r := chi.NewRouter()
    r.Use(corsMiddleware.Handler)
    r.Use(middleware.Logger)


    r.Post("/shorten", s.handleShorten)
    r.Get("/{url}", s.handleRedirect)

    return http.ListenAndServe(s.listenAddr, r)
}

func (s *Server) handleShorten(w http.ResponseWriter, r *http.Request) {
    var url types.URL

    err := json.NewDecoder(r.Body).Decode(&url)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := url.Shorten(); err != nil {
        w.Write([]byte("You need to provide an URL"))
        return
    }

    if err := s.store.InsertURL(url); err != nil {
        w.Write([]byte("Cannot expand URL"))
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    jsonData, err := json.Marshal(url)
    if err != nil{ 
        fmt.Println("Error: ", err)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write(jsonData)
}

func (s *Server) handleRedirect(w http.ResponseWriter, r *http.Request) {
    var url types.URL
    url.ShortURL = chi.URLParam(r, "url")

    s.store.GetURL(&url)
    http.Redirect(w, r, url.LongURL, http.StatusSeeOther)
}
