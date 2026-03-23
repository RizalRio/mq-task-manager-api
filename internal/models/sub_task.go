package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubTask struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TaskID      uuid.UUID      `json:"task_id" gorm:"type:uuid;not null"`
	Title       string         `json:"title" gorm:"type:varchar(150);not null"`
	Description string         `json:"description" gorm:"type:text"`
	Status      string         `json:"status" gorm:"type:varchar(50);default:'pending'"`
	Priority    string         `json:"priority" gorm:"type:varchar(20);default:'medium'"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}