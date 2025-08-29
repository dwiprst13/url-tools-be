package qr

import (
    "net/http"
    // "url-tools-be/internal/server"
    "strconv"
)


func QRHandler(w http.ResponseWriter, r *http.Request) {
    url := r.URL.Query().Get("url")
    sizeStr := r.URL.Query().Get("size")
    logoPath := "D:/Projek/Go Gank/url-tools-be/assets/logo.png"

    if url == "" {
        http.Error(w, "url parameter is required", http.StatusBadRequest)
        return
    }

    size := 256
    if sizeStr != "" {
        if s, err := strconv.Atoi(sizeStr); err == nil {
            size = s
        }
    }

    img, err := GenerateQRWithLogo(url, size, logoPath)
    if err != nil {
        http.Error(w, "failed to generate QR: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "image/png")
    w.Write(img)
}
