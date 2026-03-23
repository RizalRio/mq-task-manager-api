package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword mengenkripsi password menggunakan algoritma bcrypt
func HashPassword(password string) (string, error) {
	// Cost 10 adalah standar industri saat ini (seimbang antara keamanan dan performa)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// CheckPasswordHash membandingkan password dari user dengan hash di database
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}