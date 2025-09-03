package main

import (
	"log"
	"url-tools-be/internal/qr"
	"url-tools-be/internal/shortener"

	"github.com/gin-gonic/gin"
)

func main() {
	shortener.Init()

	router := gin.Default()
	router.Use(shortener.CORS())
	shortener.RegisterRoutes(router)
	qr.RegisterRoutes(router)
	log.Println("Server started on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

