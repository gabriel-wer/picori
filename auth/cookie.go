package auth

import (
    "net/http"
    "time"
    "crypto/rand"
    "encoding/hex"
    "errors"
)

func generateSessionID() (string, error) {
    b := make([]byte, 32) // Generate a 32-byte random string
    _, err := rand.Read(b)
    if err != nil {
        return "", err
    }
    return hex.EncodeToString(b), nil
}

func GenerateCookie() (*http.Cookie, error) {
    sessionID, err := generateSessionID()
    if err != nil {
        return &http.Cookie{}, errors.New("Error Generating a new SessionID")
    }

    cookie := &http.Cookie{
        Name: "session_cookie",
        Value: sessionID,
        Path: "/",
        Domain: "localhost",
        Expires: time.Now().Add(24 * time.Hour),
        MaxAge: 86400,
        HttpOnly: true,
        Secure: true,
    }

    return cookie, nil
}

