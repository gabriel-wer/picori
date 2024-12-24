package api

import (
	"encoding/json"
	"github.com/gabriel-wer/picori"
	errors "github.com/gabriel-wer/picori/util"
	"net/http"
)

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
