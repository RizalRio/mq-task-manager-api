package handler

import (
	"net/http"
	"task-manager-api/internal/service"
	"task-manager-api/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{service}
}

// Register godoc
// @Summary      Register User Baru
// @Description  Mendaftarkan user baru dengan email dan password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body service.RegisterRequest true "Data Registrasi"
// @Success      201  {object}  map[string]interface{} "message: Registrasi berhasil"
// @Failure      400  {object}  map[string]interface{} "error: format input salah / email sudah terdaftar"
// @Router       /auth/register [post]a
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest

	// Tangkap error validasi di sini
	if err := c.ShouldBindJSON(&req); err != nil {
		// Gunakan utilitas FormatValidationError untuk menerjemahkan error
		formattedErrors := utils.FormatValidationError(err)
		
		// Masukkan formattedErrors ke parameter "errors" pada ErrorResponse
		res := utils.ErrorResponse("Validasi input gagal", formattedErrors)
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if err := h.service.Register(req); err != nil {
		res := utils.ErrorResponse("Registrasi gagal", err.Error())
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Menggunakan SuccessResponse
	res := utils.SuccessResponse("Registrasi berhasil", nil)
	c.JSON(http.StatusCreated, res)
}

// Login godoc
// @Summary      User Login
// @Description  Autentikasi user dan mendapatkan JWT token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body service.LoginRequest true "Kredensial Login"
// @Success      200  {object}  map[string]interface{} "message: Login berhasil, token: eyJ..."
// @Failure      400  {object}  map[string]interface{} "error: format input salah"
// @Failure      401  {object}  map[string]interface{} "error: email atau password salah"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		// Implementasi terjemahan error pada fungsi Login
		formattedErrors := utils.FormatValidationError(err)
		res := utils.ErrorResponse("Validasi input gagal", formattedErrors)
		c.JSON(http.StatusBadRequest, res)
		return
	}

	token, err := h.service.Login(req)
	if err != nil {
		res := utils.ErrorResponse("Proses login gagal", err.Error())
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	// Membungkus token dalam map sebagai data
	data := map[string]string{"token": token}
	res := utils.SuccessResponse("Login berhasil", data)
	c.JSON(http.StatusOK, res)
}