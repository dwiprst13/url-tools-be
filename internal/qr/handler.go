package qr

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func QRHandler(c *gin.Context) {
	url := c.Query("url")
	sizeStr := c.Query("size")
	colorStr := c.Query("color")
	label := c.Query("label")

	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url parameter is required"})
		return
	}

	size := 256
	if sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil {
			size = s
		}
	}

	if colorStr == "" {
		colorStr = "#000000"
	}

	img, err := GenerateQRWithStyle(url, size, label, colorStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate QR: " + err.Error()})
		return
	}

	c.Data(http.StatusOK, "image/png", img)
}

