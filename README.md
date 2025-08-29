# url-tools-be

Backend API dengan Go untuk utilitas URL: **Shortener**, **Expander**, **Analyze**, dan **QR Generator**.  
Didesain modular dengan struktur `cmd/` sebagai entrypoint dan `internal/` untuk logic utama.

---

## ✨ Fitur
- 🔗 **URL Shortener**  
  - Endpoint untuk membuat URL pendek dengan alias custom (opsional).
  - Redirect otomatis dari URL pendek ke URL asli.

- 📂 **URL Expander**  
  - Mengembalikan URL asli dari sebuah short link.

- 📊 **URL Analyze**  
  - Statistik dasar: jumlah klik, sumber referer, dsb. *(coming soon)*

- 📷 **QR Generator**  
  - Generate QR code dari sebuah URL.

---

## 📂 Struktur Project

url-tools-be/
├── cmd/
│ ├── main.go # entrypoint gabungan semua fitur
│ ├── shortener/
│ │ └── main.go # entrypoint khusus shortener
│ ├── expander/
│ │ └── main.go
│ ├── analyze/
│ │ └── main.go
│ └── qr/
│ └── main.go
├── internal/
│ ├── server/ # bootstrap http.Server + middleware
│ ├── shortener/ # handler shortener
│ ├── expander/ # handler expander
│ ├── analyze/ # handler analyze
│ └── qr/ # handler qr
└── go.mod 

---

## 🚀 Menjalankan

### 1. Clone repo
```bash
git clone https://github.com/<username>/url-tools-be.git
cd url-tools-be

2. Jalankan server

Semua fitur sekaligus:

go run ./cmd/main.go


Hanya shortener:

go run ./cmd/shortener/main.go


Server default jalan di:

http://localhost:8080

📌 Endpoint
🔗 Shortener

POST /api/shorten
Request body:

{
  "url": "https://contoh.com",
  "alias": "opsional"
}


Response:

{
  "short_url": "http://localhost:8080/abc123",
  "code": "abc123"
}


GET /{code} → redirect ke URL asli.

📂 Expander

POST /api/expander
(coming soon)

📊 Analyze

GET /api/analyze/{code}
(coming soon)

📷 QR Generator

POST /api/qr
(coming soon)

⚙️ Environment Variable

PUBLIC_DOMAIN → domain publik untuk short URL (default: https://s.dprast.id).

🛠️ Tech Stack

Go 1.23+

Built-in net/http

Modular architecture dengan cmd/ & internal/

📜 License

MIT


---

📌 README ini sudah siap taruh di root project.  
Mau aku tambahin juga contoh **fetch request** (JavaScript) untuk setiap endpoint biar dev lain gampang nyoba?