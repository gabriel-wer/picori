package middleware

import (
    "net/http"
    "log"
    "time"
)

func Logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Capture the start time
        start := time.Now()

        // Call the next handler in the chain
        next.ServeHTTP(w, r)

        // Log the details after the request is processed
        log.Printf("%s %s %s %v", r.Method, r.RequestURI, r.RemoteAddr, time.Since(start))
    })
}

