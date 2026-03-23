package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	// gorm:"type:uuid;default:gen_random_uuid();primaryKey" mengatur agar DB otomatis membuat UUID
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	
	// Email harus unik dan tidak boleh kosong
	Email     string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	
	// json:"-" sangat penting! Ini mencegah password ikut terkirim saat struct ini diubah ke JSON
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	
	// Persiapan masa depan: Role (misal: "admin", "user")
	Role      string         `gorm:"type:varchar(20);default:'user'" json:"role"`
	
	// Relasi: Satu User memiliki banyak Task
	Tasks     []Task         `gorm:"foreignKey:UserID" json:"tasks,omitempty"`

	// Timestamps standar
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	
	// Soft delete dari GORM
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}