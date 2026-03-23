package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB adalah variabel global untuk menyimpan instance koneksi database
// Huruf 'D' besar menandakan variabel ini di-export dan bisa diakses dari package/folder lain
var DB *gorm.DB

// ConnectDB berfungsi untuk membuka dan memverifikasi koneksi ke PostgreSQL
func ConnectDB() {
	// 1. Memuat variabel dari file .env
	err := godotenv.Load()
	if err != nil {
		// Kita menggunakan log.Println, bukan log.Fatal di sini.
		// Alasannya: saat deploy ke server (misal Docker/AWS), kita sering tidak memakai file .env,
		// melainkan inject environment variables langsung ke OS.
		log.Println("Warning: File .env tidak ditemukan, menggunakan environment variables sistem.")
	}

	// 2. Mengambil nilai dari environment variables
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// 3. Menyusun Data Source Name (DSN) khusus untuk PostgreSQL
	// sslmode=disable digunakan untuk development lokal. Di production, idealnya ini diaktifkan (require).
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		host, user, password, dbname, port)

	// 4. Membuka koneksi menggunakan GORM
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// Jika gagal connect ke DB, aplikasi tidak punya alasan untuk terus berjalan.
		// log.Fatal akan mencetak error dan menghentikan (exit) aplikasi saat itu juga.
		log.Fatal("Gagal terhubung ke database PostgreSQL! Error: ", err)
	}

	fmt.Println("✅ Koneksi ke database PostgreSQL berhasil!")

	// 5. Menyimpan instance koneksi ke variabel global DB
	// Ini agar layer Repository nantinya bisa memanggil config.DB untuk melakukan query (CRUD).
	DB = database
}