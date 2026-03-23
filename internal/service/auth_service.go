package service

import (
	"errors"
	"task-manager-api/internal/models"
	"task-manager-api/internal/repository"
	"task-manager-api/pkg/utils"
)

// Struct untuk menerima data JSON dari request user (DTO - Data Transfer Object)
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthService interface {
	Register(req RegisterRequest) error
	Login(req LoginRequest) (string, error)
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{repo}
}

func (s *authService) Register(req RegisterRequest) error {
	// 1. Cek apakah email sudah terdaftar
	_, err := s.repo.FindByEmail(req.Email)
	if err == nil {
		return errors.New("email sudah terdaftar")
	}

	// 2. Hash password sebelum disimpan
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	// 3. Bentuk object User dan simpan
	user := &models.User{
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "user", // Default role
	}

	return s.repo.CreateUser(user)
}

func (s *authService) Login(req LoginRequest) (string, error) {
	// 1. Cari user berdasarkan email
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return "", errors.New("email atau password salah")
	}

	// 2. Cek kecocokan password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return "", errors.New("email atau password salah") // Pesan error disamakan demi keamanan agar hacker tidak tahu mana yang salah
	}

	// 3. Generate JWT Token
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}