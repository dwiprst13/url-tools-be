package server

import (
	"log"
	"net/http"
	"time"
)

type Option func(mux *http.ServeMux)

func NewServer(addr string, opts ...Option) *http.Server {
	mux := http.NewServeMux()

	for _, opt := range opts {
		opt(mux)
	}

	return &http.Server{
		Addr:         addr,
		Handler:      withCORS(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
}

func Start(srv *http.Server) {
	log.Println("listening on", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.ServeHTTP(w, r)
	})
}
