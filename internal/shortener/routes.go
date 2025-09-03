package shortener

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	api.Use(RateLimitMiddleware()) 

	api.POST("/shorten", ShortenHandler)
	router.GET("/:code", RedirectHandler)
}

