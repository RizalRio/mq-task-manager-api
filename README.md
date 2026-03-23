# MQ Task Manager API

Sistem manajemen tugas (Task Manager) skala penuh (Fullstack) mendemonstrasikan fondasi backend tingkat industri (Golang) yang mengkombinasikan *Data Isolation*, sistem otorisasi berbasis peran (RBAC), dan eksekusi latar belakang asinkron, yang dipadukan dengan pemrosesan AI. Repositori ini difokuskan pada **Backend Layer (Golang API)** dari sistem Task Manager.

## 🚀 Fitur Utama

- **Autentikasi & Otorisasi:** JWT & *Role-Based Access Control* (RBAC).
- **Manajemen Tugas:** CRUD tugas dengan *Pagination*, *Filtering*, dan *Soft Delete*.
- **Kecerdasan Buatan (AI):** Terintegrasi langsung dengan Google Generative AI (Gemini) untuk pemecahan tugas (*Sub-Tasks*).
- **Kolaborasi Tim:** *Task Sharing* yang diterapkan melalui skema *Pivot Table*.
- **Background Jobs:** Penjadwalan latar belakang (Cron Jobs) secara terotomatisasi untuk meninjau penanda tugas berstatus *overdue*.
- **Upload File:** Sistem unggahan (*file attachments*) yang dapat dinikmati demi kelengkapan tugas.

## 🛠️ Teknologi & Tools

- **Bahasa Pemrograman:** Golang 1.25+
- **Framework REST API:** Gin HTTP Framework (`gin-gonic/gin`)
- **Database:** PostgreSQL dengan ORM (`gorm.io/gorm`)
- **Migrasi Database:** `golang-migrate` dijalankan melalui `Makefile`
- **Keamanan:** Bcrypt, JWT (`golang-jwt/jwt/v5`)
- **Background Job:** Robfig Cron v3 (`robfig/cron/v3`)
- **AI Teks:** Google Generative AI SDK (`gemini-1.5-flash-latest`)
- **Dokumentasi API:** Swagger (`swaggo/swag`, `swaggo/gin-swagger`)

## 📂 Struktur Direktori (Layered Architecture)

Arsitektur aplikasi ini menggunakan *Clean / Layered Architecture*:
- `cmd/api/`: *Entry point*, manajemen *Dependency Injection*, *routing*, dan inisialisasi *worker*.
- `internal/repository/`: Lapisan *Data Access*, eksekusi *query* PostgreSQL, agregasi data.
- `internal/service/`: Lapisan Logika Bisnis (*Business Logic*), pelindung *guard clause* RBAC, injeksi AI *prompt engineering*.
- `internal/handler/`: Lapisan Pengendali HTTP, *parsing request*, standarisasi respons.
- `internal/middleware/`: Interseptor validasi token JWT dan pengaturan CORS.
- `internal/worker/`: Lapisan eksekusi *Background Task* (Cron Jobs).
- `pkg/`: Kode utilitas pendukung (seperti fungsi global).

## ⚙️ Persyaratan Sistem

- [Go](https://golang.org/dl/) (versi 1.25 atau yang lebih terbaru)
- [PostgreSQL](https://www.postgresql.org/download/) terinstal
- [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) CLI (untuk menjalanan migrasi *database*)
- [Make](https://www.gnu.org/software/make/) utilitas untuk menjalankan *Makefile*

## 📦 Panduan Instalasi & Startup Lokal

1. **Clone repositori ini** dan masuk ke direktori proyek lokal:
   ```bash
   git clone <repo-url>
   cd mq-task-manager-api
   ```

2. **Konfigurasi Environment Variables**
   Salin *template environment variables* atau buat `.env` baru pada *root* proyek. Konfigurasi wajib memuat blok berikut (sesuaikan dengan _database local_ Anda):
   ```env
   # Database Configuration
   DB_HOST=localhost
   DB_USER=postgres
   DB_PASSWORD=password
   DB_NAME=gk_task_manager
   DB_PORT=5432

   # Application Configuration
   PORT=8080
   JWT_SECRET=rahasia_super_aman_untuk_token_jwt

   # Gemini API Key (Untuk fitur AI)
   GEMINI_API_KEY=AI...
   ```

3. **Inisialisasi Database**
   Pastikan Anda sudah membuat *database* `gk_task_manager` di PostgreSQL Anda terlebih dahulu. Setalah itu, jalankan migrasi yang telah dirangkai di dalam *Makefile*:
   ```bash
   make migrate-up
   ```

4. **Unduh Dependensi Golang**
   ```bash
   go mod tidy
   ```

5. **Jalankan Aplikasi Backend API**
   ```bash
   go run cmd/api/main.go
   ```
   
   *Server middleware* REST API akan berjalan pada tautan `http://localhost:8080`.

## 📌 Catatan Referensi Lebih Lengkap

Jika Anda ingin mengetahui detail konotasi di balik *front-end layer* (Next.js) dan arsitektur eksekusi kolaborasi tugas secara menyeluruh, Anda bisa merujuk ke dalam doks spesifikasi pada berkas [project_context.md](./project_context.md).
