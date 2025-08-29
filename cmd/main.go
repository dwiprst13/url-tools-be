package main

import (
	"url-tools-be/internal/server"
	"url-tools-be/internal/shortener"
	"url-tools-be/internal/qr"
)

func main() {
	shortener.Init()

	srv := server.NewServer(":8080",
		shortener.Routes(),
		qr.Routes(),
	)

	server.Start(srv)
}
