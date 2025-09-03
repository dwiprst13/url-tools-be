package shortener

import (
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type shortenReq struct {
	URL   string `json:"url"`
	Alias string `json:"alias,omitempty"`
}

type shortenResp struct {
	ShortURL string `json:"short_url"`
	Code     string `json:"code"`
}

var (
	store        Storage
	limiter      *RateLimiter
	codePattern  = regexp.MustCompile(`^[0-9A-Za-z_-]{3,32}$`)
	codeLength   = 6
	publicDomain string
)

func Init() {
	store = NewMemoryStore()
	limiter = NewRateLimiter(5, 10, time.Second)

	publicDomain = os.Getenv("PUBLIC_DOMAIN")
	if publicDomain == "" {
		publicDomain = "https://s.dprast.id"
	}
}

func ShortenHandler(c *gin.Context) {
	var req shortenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	longURL, err := NormalizeURL(strings.TrimSpace(req.URL))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
		return
	}

	if existCode, ok := store.FindByURL(longURL); ok && req.Alias == "" {
		resp := shortenResp{
			ShortURL: publicDomain + "/" + existCode,
			Code:     existCode,
		}
		c.JSON(http.StatusOK, resp)
		return
	}

	code := strings.TrimSpace(req.Alias)
	if code != "" {
		if !codePattern.MatchString(code) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid alias (3-32 alnum/_-)})"})
			return
		}
		if err := store.Save(code, longURL); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "alias already in use"})
			return
		}
	} else {
		for {
			gen, err := GenerateCode(codeLength)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "generator error"})
				return
			}
			if err := store.Save(gen, longURL); err == nil {
				code = gen
				break
			}
		}
	}
	resp := shortenResp{
		ShortURL: publicDomain + "/" + code,
		Code:     code,
	}
	c.JSON(http.StatusCreated, resp)
}

func RedirectHandler(c *gin.Context) {
	code := c.Param("code") 

	if code == "" || code == "favicon.ico" || strings.HasPrefix(code, "api/") {
		c.String(http.StatusOK, "URL Shortener OK")
		return
	}

	longURL, err := store.Get(code)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.Redirect(http.StatusMovedPermanently, longURL)
}

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow(c.Request) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}
		c.Next() 
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

