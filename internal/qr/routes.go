package qr

import (
	"net/http"
	"url-tools-be/internal/server"
)

func Routes() server.Option {
	return func(mux *http.ServeMux) {
		mux.HandleFunc("/qr", QRHandler)
	}
}