package api

import (
	"encoding/json"
	"net/http"
	"time"
)

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

	ldap, err := NewLDAPConnection()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = ldap.AuthenticateUser(login.Username, login.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	sessionID, err := GenerateSessionID()
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

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}
