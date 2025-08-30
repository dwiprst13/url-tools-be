package qr

import (
    "net/http"
    "strconv"
)

func QRHandler(w http.ResponseWriter, r *http.Request) {
    url := r.URL.Query().Get("url")
    sizeStr := r.URL.Query().Get("size")
    colorStr := r.URL.Query().Get("color")
    label := r.URL.Query().Get("label")

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

    if colorStr == "" {
        colorStr = "#000000" 
    }

    img, err := GenerateQRWithStyle(url, size, label, colorStr)
    if err != nil {
        http.Error(w, "failed to generate QR: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "image/png")
    w.Write(img)
}
