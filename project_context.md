# Project Context: Task Manager Fullstack System

## 1. Deskripsi Project

Sistem manajemen tugas (Task Manager) skala penuh (Fullstack) dengan fitur autentikasi pengguna, kolaborasi tim, dan asisten kecerdasan buatan. Sistem ini mendemonstrasikan fondasi _backend_ tingkat industri (Golang) yang mengkombinasikan _Data Isolation_, sistem otorisasi berbasis peran (RBAC), dan eksekusi latar belakang asinkron, yang dipadukan dengan antarmuka _frontend_ modern (Next.js) yang interaktif, responsif, dan menerapkan pola _Client-Side Rendering_ (CSR) serta Route Protection.

## 2. Tech Stack Utama

**Backend (RESTful API):**

- **Bahasa Pemrograman:** Golang
- **Web Framework:** Gin HTTP Framework (`github.com/gin-gonic/gin`)
- **Database:** PostgreSQL dengan GORM (`gorm.io/gorm`)
- **Database Migration Tool:** `golang-migrate`
- **Keamanan:** Bcrypt & JWT (`golang-jwt/jwt/v5`)
- **Kecerdasan Buatan (AI):** Google Generative AI SDK (`gemini-1.5-flash-latest`)
- **Background Job:** Robfig Cron v3 (`github.com/robfig/cron/v3`)

**Frontend (Client Interface):**

- **Framework:** Next.js (App Router) dengan React
- **Styling:** Tailwind CSS
- **HTTP Client:** Axios (dengan Interceptors untuk injeksi JWT)
- **State Management:** React Hooks (`useState`, `useEffect`)
- **Ikon & Tipografi:** Lucide React

## 3. Arsitektur Sistem (Layered & Client-Server Architecture)

Project ini memisahkan logika bisnis dari antarmuka presentasi:

**A. Backend Layer (Golang):**

- `cmd/api/`: _Entry point_, _Dependency Injection_, _routing_, inisialisasi _worker_.
- `internal/repository/`: Lapisan _Data Access_, _query_ PostgreSQL, agregasi data.
- `internal/service/`: Lapisan Logika Bisnis, _guard clause_ RBAC, _prompt engineering_ AI.
- `internal/handler/`: Lapisan Pengendali HTTP, _parsing request_, respons terstandarisasi.
- `internal/middleware/`: Interseptor validasi token JWT dan CORS.
- `internal/worker/`: Lapisan _Background Task_ (Cron Jobs).

**B. Frontend Layer (Next.js):**

- `src/app/auth/`: Modul rute publik untuk proses Registrasi dan Login.
- `src/app/dashboard/`: Modul rute terproteksi untuk Dasbor Utama (Daftar Tugas) dan Detail Tugas dinamis (`[id]`).
- `src/lib/axios.ts`: Jantung komunikasi HTTP yang menempelkan token sesi ke setiap _request_ secara otomatis.
- `src/middleware.ts`: Penjaga gerbang (_Route Guard_) berbasis _server_ yang mencegah akses tidak sah ke area dasbor.

## 4. Status Implementasi Fitur

**Sisi Backend (API):**

- [x] Autentikasi JWT & Otorisasi RBAC.
- [x] CRUD Tugas dengan Pagination, Filtering, dan _Soft Delete_.
- [x] Integrasi AI Assistant untuk pemecahan tugas (Sub-Tasks).
- [x] Sistem Unggah Lampiran (File Attachment).
- [x] Kolaborasi Tim (Task Sharing / Pivot Table).
- [x] Penjadwalan Latar Belakang (Cron Jobs) untuk status _overdue_.

**Sisi Frontend (UI/UX):**

- [x] _Setup_ Next.js & Tailwind CSS.
- [x] Halaman Autentikasi (Login & Register) dengan _Error Handling_.
- [x] _Middleware_ Proteksi Rute berbasis Cookies.
- [x] Dasbor Utama: Menampilkan Daftar Tugas dengan lencana prioritas dinamis.
- [x] Modal Interaktif: Formulir Pembuatan dan Pembaruan (Edit) Tugas.
- [x] Halaman Detail Tugas (Dynamic Routing `[id]`).
- [x] Antarmuka Unggah Lampiran dan Pemanggilan AI Breakdown.
- [x] Tampilan Pil Badge Kolaborator & Modal Undangan Anggota Tim.
- [x] Fungsi Interaktif Tandai Selesai dan Hapus Tugas Permanen.

## 5. Rencana Pengembangan Lanjutan (Roadmap)

- [ ] **Sistem Pendukung Keputusan (DSS) Prioritas Tugas:** Mengimplementasikan algoritma analitik (AHP dan TOPSIS) di _backend_ untuk memberikan rekomendasi tugas paling kritis, ditampilkan di dasbor _frontend_.
- [ ] **Containerization (Docker):** Membungkus layanan Golang, database PostgreSQL, dan _frontend_ Next.js ke dalam kontainer menggunakan `docker-compose` agar siap tayang (_deploy_).
- [ ] **Notifikasi Real-time (WebSockets):** Menerapkan komunikasi _Full-Duplex_ agar kolaborator mendapatkan pembaruan data secara instan tanpa perlu memuat ulang halaman.
