package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type longURL struct { 
    URL string `json:"url"`
}

type shortURL struct { 
    ShortURL string `json:"shortURL"`
    LongURL string `json:"longURL"`


}

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
    var longURL longURL

    err := json.NewDecoder(r.Body).Decode(&longURL)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    shortcode := shortener.shortenURL(longURL.URL)
    short := fmt.Sprintf("http://localhost:3001/%s", shortcode)

    shortenedUrl := shortURL{
        ShortURL: short,
        LongURL: longURL.URL,
    }
    jsonData, err := json.Marshal(shortenedUrl)
    if err != nil{ 
        fmt.Println("Error: ", err)
        return
    }

    w.Write(jsonData)
    
}

func expand(w http.ResponseWriter, r *http.Request) {
    shorturl := chi.URLParam(r, "url")

    longcode := shortener.expandURL(shorturl)
    w.Write([]byte(longcode))
    
}

func main() {
    r := chi.NewRouter()
    r.Use(middleware.Logger)

    r.Post("/shorten", shorten)
    r.Get("/{url}", expand)

    http.ListenAndServe(":3001", r)
}
