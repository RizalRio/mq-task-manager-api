package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// RequireAuth adalah fungsi middleware untuk memproteksi endpoint
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Mengambil header Authorization dari request HTTP
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak: Header Authorization tidak ditemukan"})
			c.Abort() // Menghentikan request agar tidak lanjut ke handler berikutnya
			return
		}

		// 2. Memastikan formatnya adalah "Bearer <token>"
		// strings.SplitN memecah string menjadi 2 bagian berdasarkan spasi
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Format token tidak valid. Gunakan format: Bearer <token>"})
			c.Abort()
			return
		}

		// Mengisolasi string token
		tokenString := parts[1]
		secret := os.Getenv("JWT_SECRET")

		// 3. Melakukan parsing dan validasi algoritma token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Memastikan token ditandatangani menggunakan algoritma HMAC (HS256)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		// 4. Menangani error jika token kedaluwarsa (expired) atau rusak
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau sudah kedaluwarsa"})
			c.Abort()
			return
		}

		// 5. Mengekstrak payload/claims dari dalam token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Gagal memproses data di dalam token"})
			c.Abort()
			return
		}

		// 6. Menyisipkan user_id dan role ke dalam Gin Context
		// Ini sangat penting agar Handler Task nantinya tahu siapa user yang sedang melakukan request
		// tanpa harus mengirimkan ID di body request.
		c.Set("user_id", claims["user_id"])
		c.Set("role", claims["role"])

		// 7. Melanjutkan request ke Handler tujuan
		c.Next()
	}
}