package worker

import (
	"log"
	"time"

	"task-manager-api/internal/repository"

	"github.com/robfig/cron/v3"
)

// StartTaskWorker menginisialisasi dan menjalankan jadwal latar belakang
func StartTaskWorker(taskRepo repository.TaskRepository) {
	// Set zona waktu presisi
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.UTC
		log.Println("[CRON WARNING] Gagal memuat zona waktu Asia/Jakarta, menggunakan UTC")
	}

	// Membuat instance cron baru dengan zona waktu yang sudah diset
	c := cron.New(cron.WithLocation(loc))

	// Mengatur jadwal eksekusi (Format: Menit Jam Tanggal Bulan Hari)
	// "* * * * *" artinya dieksekusi SETIAP 1 MENIT (cocok untuk uji coba)
	_, err = c.AddFunc("* * * * *", func() {
		// Panggil fungsi Bulk Update dari Repository
		rowsAffected, err := taskRepo.UpdateOverdueTasks()
		if err != nil {
			log.Printf("[CRON ERROR] Gagal mengecek tenggat waktu tugas: %v\n", err)
			return
		}
		
		// Hanya tampilkan log jika ada data yang benar-benar diubah
		if rowsAffected > 0 {
			log.Printf("[CRON INFO] Berhasil menandai %d tugas menjadi OVERDUE\n", rowsAffected)
		}
	})

	if err != nil {
		log.Fatalf("Gagal mendaftarkan fungsi Cron Job: %v", err)
	}

	// Mulai jalankan cron di background (asynchronous)
	c.Start()
	log.Println("[CRON INFO] Task Worker Scheduler berjalan di background...")
}