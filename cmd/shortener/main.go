package main

import (
	"url-tools-be/internal/server"
	"url-tools-be/internal/shortener"
)

func main() {
	srv := server.NewServer(":8080",
		shortener.Routes(),
	)
	server.Start(srv)
}
