package server

import (
	"encoding/json"
	"net/http"

	"github.com/gabriel-wer/picori/storage"
	"github.com/gabriel-wer/picori/picori"
    "github.com/gabriel-wer/picori/middleware"
    "github.com/gabriel-wer/picori/auth"
	errors "github.com/gabriel-wer/picori/util"
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
    mux := http.NewServeMux()

    mux.HandleFunc("POST /shorten", s.handleShorten)
    mux.HandleFunc("POST /expand", s.handleExpand)
    mux.HandleFunc("GET /{url}", s.handleRedirect)
    mux.HandleFunc("POST /login", s.handleLogin)
    
    midMux := middleware.Chain(mux, middleware.Logging,
                                    middleware.CORS)

    return http.ListenAndServe(s.listenAddr, midMux)
}
type LoginForm struct {
    Username string `json:"username"`
    Password string `json:"password"`
}
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
    var login LoginForm
    err := json.NewDecoder(r.Body).Decode(&login)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if login.Username == "" || login.Password == "" {
        http.Error(w, "Missing username or password", http.StatusBadRequest)
        return
    }

    ldap, err := auth.NewLDAPConnection()
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    err = ldap.AuthenticateUser(login.Username, login.Password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    cookie, err := auth.GenerateCookie()
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    w.WriteHeader(http.StatusOK)
    http.SetCookie(w, cookie)
    w.Write([]byte(cookie.Value))
}

func (s *Server) handleShorten(w http.ResponseWriter, r *http.Request) {
    var url picori.URL

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
        w.Write([]byte("Cannot Shorten URL"))
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write(errors.JsonData(w, url))
}

func (s *Server) handleExpand(w http.ResponseWriter, r *http.Request) {
    var url picori.URL

    err := json.NewDecoder(r.Body).Decode(&url)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := s.store.GetURL(&url); err != nil {
        w.Write([]byte("Cannot expand URL"))
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    jsonData, err := json.Marshal(url)
    if err != nil { 
        w.Write([]byte("Cannot marshall JSON Data"))
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }

    w.WriteHeader(http.StatusOK)
    w.Write(jsonData)
}

func (s *Server) handleRedirect(w http.ResponseWriter, r *http.Request) {
    var url picori.URL
    url.ShortURL = r.PathValue("url")
    //TODO: Fix redirects
    s.store.GetURL(&url)
    w.WriteHeader(http.StatusPermanentRedirect)
    http.Redirect(w, r, url.LongURL, 301)
}
