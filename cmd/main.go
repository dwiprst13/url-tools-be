package main

import (
	"log"
	"url-tools-be/internal/auth"
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
	api := router.Group("/api")
	api.Use(shortener.RateLimitMiddleware())
	api.POST("/shorten", shortener.ShortenHandler)
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/login", auth.LoginHandler)
		authGroup.GET("/profile", auth.AuthMiddleware(), auth.ProfileHandler)
	}
	log.Println("Server started on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
	router.GET("/:code", shortener.RedirectHandler)
	qr.RegisterRoutes(router)
}
