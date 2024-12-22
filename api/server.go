package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gabriel-wer/picori"
	"github.com/gabriel-wer/picori/api/auth"
	"github.com/gabriel-wer/picori/api/middleware"
	"github.com/gabriel-wer/picori/storage"
	errors "github.com/gabriel-wer/picori/util"
)

type Server struct {
	listenAddr string
	store      *storage.Sqlite
}

func NewServer(listenAddr string, store *storage.Sqlite) *Server {
	return &Server{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /v1/shorten", s.handleShorten)
	mux.HandleFunc("POST /v1/expand", s.handleExpand)
	mux.HandleFunc("GET /v1/{url}", s.handleRedirect)
	mux.HandleFunc("POST /v1/login", s.handleLogin)
	mux.HandleFunc("GET /v1/list", middleware.Authentication(s.handleList, s.store))

	midMux := middleware.Chain(mux, middleware.CORS)
	midMux = middleware.Logging(midMux)

	go s.serveStatic()
	return http.ListenAndServe(s.listenAddr, midMux)
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

func (s *Server) handleList(w http.ResponseWriter, r *http.Request) {
	urls, err := s.store.ListURL()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(urls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	type LoginForm struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
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

	sessionID, err := auth.GenerateSessionID()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = s.store.SaveCookie(login.Username, sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	cookie := &http.Cookie{
		Name:     "cookayyy",
		Value:    sessionID,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		MaxAge:   86400,
		SameSite: http.SameSiteLaxMode,
		Domain:   "localhost",
	}

	resp := fmt.Sprintf(`{"cookie": "%s" }`, cookie.Value)
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resp))
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
