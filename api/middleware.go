package api

import (
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/gabriel-wer/picori/storage"
)

type middleware func(http.Handler) http.Handler

type router struct {
	*http.ServeMux
	chain []middleware
}

func Chain(handler http.Handler, middlewares ...middleware) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}

	return handler
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Printf("%s %s %s %v", r.Method, r.RequestURI, r.RemoteAddr, time.Since(start))
	})
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTION")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Authentication(next http.HandlerFunc, store *storage.Sqlite) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("cookayyy")
		if err != nil {
			http.Error(w, "Unauthorized1", http.StatusUnauthorized)
			return
		}

		err = store.CheckCookie(cookie.Value)
		if err != nil {
			http.Error(w, "Unauthorized2", http.StatusUnauthorized)
			return
		}

		next(w, r)
	})
}

func NewRouter(mx ...middleware) *router {
	return &router{ServeMux: &http.ServeMux{}, chain: mx}
}

func (r *router) Use(mx ...middleware) {
	r.chain = append(r.chain, mx...)
}

func (r *router) Group(fn func(r *router)) {
	fn(&router{ServeMux: r.ServeMux, chain: slices.Clone(r.chain)})
}

func (r *router) Get(path string, fn http.HandlerFunc, mx ...middleware) {
	r.handle(http.MethodGet, path, fn, mx)
}

func (r *router) Post(path string, fn http.HandlerFunc, mx ...middleware) {
	r.handle(http.MethodPost, path, fn, mx)
}

func (r *router) Put(path string, fn http.HandlerFunc, mx ...middleware) {
	r.handle(http.MethodPut, path, fn, mx)
}

func (r *router) Delete(path string, fn http.HandlerFunc, mx ...middleware) {
	r.handle(http.MethodDelete, path, fn, mx)
}

func (r *router) Head(path string, fn http.HandlerFunc, mx ...middleware) {
	r.handle(http.MethodHead, path, fn, mx)
}

func (r *router) Options(path string, fn http.HandlerFunc, mx ...middleware) {
	r.handle(http.MethodOptions, path, fn, mx)
}

func (r *router) handle(method, path string, fn http.HandlerFunc, mx []middleware) {
	r.Handle(method+" "+path, r.wrap(fn, mx))
}

func (r *router) wrap(fn http.HandlerFunc, mx []middleware) (out http.Handler) {
	out, mx = http.Handler(fn), append(slices.Clone(r.chain), mx...)

	slices.Reverse(mx)

	for _, m := range mx {
		out = m(out)
	}

	return
}
