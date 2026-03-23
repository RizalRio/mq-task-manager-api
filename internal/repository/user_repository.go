package repository

import (
	"task-manager-api/internal/models"

	"gorm.io/gorm"
)

// Gunakan Interface agar future-proof (mudah untuk Unit Testing)
type UserRepository interface {
	CreateUser(user *models.User) error
	FindByEmail(email string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository adalah constructor untuk repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

// CreateUser menyimpan user baru ke database
func (r *userRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

// FindByEmail mencari user berdasarkan email (dibutuhkan saat login dan validasi register)
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}