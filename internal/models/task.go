package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	
	// Foreign key yang mengarah ke tabel users
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	
	Title       string         `gorm:"type:varchar(150);not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	
	// Status tugas (pending, in_progress, completed)
	Status      string         `gorm:"type:varchar(20);default:'pending'" json:"status"`
	
	// Persiapan masa depan: Prioritas (low, medium, high) dan Tenggat Waktu
	Priority    string         `gorm:"type:varchar(20);default:'medium'" json:"priority"`
	AttachmentURL string         `json:"attachment_url" gorm:"type:varchar(255)"`
	
	// Relasi: Satu Task memiliki banyak SubTask
	SubTasks	[]SubTask `json:"sub_tasks" gorm:"foreignKey:TaskID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	DueDate     *time.Time     `json:"due_date,omitempty"` // Menggunakan pointer (*) agar bisa bernilai nil/null

	Collaborators []User         `json:"collaborators" gorm:"many2many:task_collaborators;"`

	CreatedAt   time.Time      `json:"created_at"`	
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}