// internal/shortener/routes.go
package shortener

import (
	"net/http"
	"url-tools-be/internal/server"
)

func Routes() server.Option {
	return func(mux *http.ServeMux) {
		mux.HandleFunc("/api/shorten", rateLimit(jsonOnly(shortenHandler)))
		mux.HandleFunc("/", redirectHandler)
	}
}
