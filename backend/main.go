package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

var db *sql.DB;

type URL struct { 
    ShortURL string `json:"shorturl" gorm:"primaryKey;index"`
    LongURL string `json:"longurl"`
}

func shorten(w http.ResponseWriter, r *http.Request) {
    var url URL

    err := json.NewDecoder(r.Body).Decode(&url)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    hash := sha1.Sum([]byte(url.LongURL))
    url.ShortURL = fmt.Sprintf("%x", hash)[:6]

    _, err = db.Exec("INSERT INTO url (shorturl, longurl) VALUES ($1, $2)", url.ShortURL, url.LongURL)
    if err != nil {
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

func redirect(w http.ResponseWriter, r *http.Request) {
    var url URL
    url.ShortURL = chi.URLParam(r, "url")

    db.QueryRow("SELECT longurl FROM url WHERE shorturl = $1", url.ShortURL).Scan(&url.LongURL)
    if !strings.Contains(url.LongURL, "https") || !strings.Contains(url.LongURL, "http") {
        url.LongURL = "https://" + url.LongURL
    }
    fmt.Println(url.LongURL)
    http.Redirect(w, r, url.LongURL, http.StatusSeeOther)
}

func setup() *sql.DB{
    dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=picori sslmode=disable", 
                        os.Getenv("DB_HOST"), 
                        os.Getenv("DB_USER"), 
                        os.Getenv("DB_PASSWORD"), 
                        os.Getenv("DB_PORT"))
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        panic(err)
    }

    _, err = db.Exec("CREATE TABLE IF NOT EXISTS url (shorturl TEXT, longurl TEXT, PRIMARY KEY(shorturl))")
    if err != nil {
        panic(err)
    }

    return db
}

func main() {
    db = setup()
    defer db.Close()

    r := chi.NewRouter()
    r.Use(middleware.Logger)

    r.Post("/shorten", shorten)
    r.Get("/{url}", redirect)

    port := fmt.Sprintf(":%s", os.Getenv("PICORI_PORT"))
    http.ListenAndServe(port, r)
}
