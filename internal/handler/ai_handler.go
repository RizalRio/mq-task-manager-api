package handler

import (
	"net/http"
	"task-manager-api/internal/service"
	"task-manager-api/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AIHandler struct {
	aiService service.AIService
}

func NewAIHandler(aiService service.AIService) *AIHandler {
	return &AIHandler{aiService}
}

// GenerateTaskBreakdown godoc
// @Summary      Generate Sub-Tasks dengan AI
// @Description  Memecah tugas besar menjadi langkah-langkah kecil menggunakan Google Gemini AI
// @Tags         AI Assistant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body service.AITaskBreakdownRequest true "Judul Tugas Besar"
// @Success      200  {object}  utils.DefaultResponse
// @Failure      400  {object}  utils.DefaultResponse
// @Failure      401  {object}  utils.DefaultResponse
// @Failure      500  {object}  utils.DefaultResponse
// @Router       /ai/generate-tasks [post]
func (h *AIHandler) GenerateTaskBreakdown(c *gin.Context) {
	// Memastikan hanya user yang terautentikasi yang bisa mengakses fitur AI
	_, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Autentikasi gagal", "User tidak valid"))
		return
	}

	var req service.AITaskBreakdownRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		formattedErrors := utils.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Validasi input gagal", formattedErrors))
		return
	}

	// Memanggil AI Service. Kita meneruskan c.Request.Context() agar context HTTP 
	// tetap terjaga saat menghubungi server Google.
	suggestions, err := h.aiService.GenerateTaskBreakdown(c.Request.Context(), req.TaskTitle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Gagal memproses permintaan AI", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Berhasil mendapatkan saran tugas dari AI", suggestions))
}