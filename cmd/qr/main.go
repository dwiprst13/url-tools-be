package main

import (
	"log"
	"net/http"

	"url-tools-be/internal/qr"
)

func main() {
	logoPath := "assets/logo.svg" 

	http.HandleFunc("/qr", qr.QRHandler(logoPath))

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
