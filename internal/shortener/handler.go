package shortener

import (
	"encoding/json"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type shortenReq struct {
	URL   string `json:"url`
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

func Routes() func(mux *http.ServeMux) {
	return func(mux *http.ServeMux) {
		mux.HandleFunc("/api/shorten", rateLimit(jsonOnly(shortenHandler)))
		mux.HandleFunc("/", redirectHandler)
	}
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, http.StatusMethodNotAllowed, "methode not alowwed")
		return
	}

	var req shortenReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid json")
		return
	}

	longURL, err := NormalizeURL(strings.TrimSpace(req.URL))
	if err != nil {
		httpError(w, http.StatusBadRequest, "invalid url")
		return
	}

	if existCode, ok := store.FindByURL(longURL); ok && req.Alias == "" {
		resp := shortenResp{
			ShortURL: publicDomain + "/" + existCode,
			Code:     existCode,
		}
		writeJSON(w, http.StatusOK, resp)
		return
	}

	code := strings.TrimSpace(req.Alias)
	if code != "" {
		if !codePattern.MatchString(code) {
			httpError(w, http.StatusBadRequest, "invalid alias (3-32 alnum/_-)")
			return
		}
		if err := store.Save(code, longURL); err != nil {
			httpError(w, http.StatusBadRequest, "alias already in use")
			return
		}
	} else {
		for {
			gen, err := GenerateCode(codeLength)
			if err != nil {
				httpError(w, http.StatusBadRequest, "generator error")
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
	writeJSON(w, http.StatusCreated, resp)
}

func redirectHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet {
		httpError(w, http.StatusMethodNotAllowed, "methode not allowed")
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" || path == "favicon.ico" || strings.HasPrefix(path, "api/"){
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("url shortener ok"))
		return
	}
	code := path
	longURL, err := store.Get(code)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, longURL, http.StatusMovedPermanently)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func httpError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func rateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow(r) {
			httpError(w, http.StatusTooManyRequests, "rate limit exceeded")
			return
		}
		next.ServeHTTP(w, r)
	}
}

func jsonOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
			httpError(w, http.StatusUnsupportedMediaType, "content-type must be application/json")
			return
		}
		next.ServeHTTP(w, r)
	}
}

func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.ServeHTTP(w, r)
	})
}