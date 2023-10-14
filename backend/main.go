package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

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

    result := db.Create(&url)
    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
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

    db.Where("short_url = ?", url.ShortURL).Find(&url)
    http.Redirect(w, r, url.LongURL, http.StatusSeeOther)
}

func main() {
    dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s", 
                        os.Getenv("PICORI_HOST"), 
                        os.Getenv("PICORI_USER"), 
                        os.Getenv("PICORI_PASSWORD"), 
                        os.Getenv("PICORI_PORT"))
    var err error
    db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

    if err != nil {
        panic(err)
    }

    db.AutoMigrate(&URL{})

    r := chi.NewRouter()
    r.Use(middleware.Logger)

    r.Post("/shorten", shorten)
    r.Get("/{url}", redirect)

    http.ListenAndServe(":3001", r)
}
