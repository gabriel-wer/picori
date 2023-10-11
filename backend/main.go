package main

import (
	"fmt"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
    "crypto/sha1"
)

type LinkShortener struct { 
    urlMap map[string]string
}

func NewLinkShortener() *LinkShortener {
    return &LinkShortener{urlMap: make(map[string]string)}
}

var shortener = NewLinkShortener()

func (ls *LinkShortener) shortenURL(longURL string) string{
    hash := sha1.Sum([]byte(longURL))
    shortCode := fmt.Sprintf("%x", hash)[:6]

    ls.urlMap[shortCode] = longURL

    return shortCode
}

func (ls *LinkShortener) expandURL(shortURL string) string{
    longURL := ls.urlMap[shortURL]
    return longURL
}

func shorten(w http.ResponseWriter, r *http.Request) {
    longURL := chi.URLParam(r, "url")

    shortcode := shortener.shortenURL(longURL)
    shortURL := fmt.Sprintf("http://localhost:3001/%s", shortcode)
    w.Write([]byte(shortURL))
    
}

func expand(w http.ResponseWriter, r *http.Request) {
    shorturl := chi.URLParam(r, "url")

    longcode := shortener.expandURL(shorturl)
    w.Write([]byte(longcode))
    
}

func main() {
    r := chi.NewRouter()
    r.Use(middleware.Logger)


    r.Route("/", func(r chi.Router){
        r.Post("/{url}", shorten)
        r.Get("/{url}", expand)
    })

    http.ListenAndServe(":3001", r)
}
