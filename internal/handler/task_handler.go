package handler

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"task-manager-api/internal/service"
	"task-manager-api/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TaskHandler struct {
	taskService service.TaskService
}

func NewTaskHandler(taskService service.TaskService) *TaskHandler {
	return &TaskHandler{taskService}
}

func getUserID(c *gin.Context) (uuid.UUID, error) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, errors.New("user_id not found in context")
	}
	return uuid.Parse(userIDStr.(string))
}

// CreateTask godoc
// @Summary      Buat Tugas Baru
// @Description  Menambahkan tugas baru untuk user yang sedang login
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body service.CreateTaskRequest true "Data Tugas"
// @Success      201  {object}  utils.DefaultResponse
// @Failure      400  {object}  utils.DefaultResponse
// @Failure      401  {object}  utils.DefaultResponse
// @Router       /tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Autentikasi gagal", "User tidak valid atau token rusak"))
		return
	}

	var req service.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Menggunakan penerjemah validasi
		formattedErrors := utils.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Validasi input gagal", formattedErrors))
		return
	}

	task, err := h.taskService.CreateTask(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Gagal membuat tugas", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("Tugas berhasil dibuat", task))
}

// GetTasks godoc
// @Summary      Ambil Daftar Tugas
// @Description  Mengambil daftar tugas milik user dengan dukungan pagination dan filtering
// @Tags         Tasks
// @Produce      json
// @Security     BearerAuth
// @Param        page     query     int     false  "Nomor Halaman (Default: 1)"
// @Param        limit    query     int     false  "Jumlah Data per Halaman (Default: 10)"
// @Param        status   query     string  false  "Filter berdasarkan status (pending, in_progress, completed)"
// @Param        priority query     string  false  "Filter berdasarkan prioritas (low, medium, high)"
// @Success      200      {object}  utils.DefaultResponse
// @Failure      401      {object}  utils.DefaultResponse
// @Router       /tasks [get]
func (h *TaskHandler) GetTasks(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Autentikasi gagal", "User tidak valid"))
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	status := c.Query("status")
	priority := c.Query("priority")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	filterParams := service.TaskFilterParams{
		Page:     page,
		Limit:    limit,
		Status:   status,
		Priority: priority,
	}

	result, err := h.taskService.GetTasksByUser(userID, filterParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Gagal mengambil daftar tugas", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Berhasil mengambil daftar tugas", result))
}

// GetTaskByID godoc
// @Summary      Ambil Detail Tugas
// @Description  Mengambil spesifik satu tugas berdasarkan ID
// @Tags         Tasks
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "UUID Tugas"
// @Success      200  {object}  utils.DefaultResponse
// @Failure      400  {object}  utils.DefaultResponse
// @Failure      401  {object}  utils.DefaultResponse
// @Failure      404  {object}  utils.DefaultResponse
// @Router       /tasks/{id} [get]
func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Autentikasi gagal", "User tidak valid"))
		return
	}

	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Format ID tugas tidak valid", err.Error()))
		return
	}

	task, err := h.taskService.GetTaskByID(taskID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Tugas tidak ditemukan", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Berhasil mengambil detail tugas", task))
}

// UpdateTask godoc
// @Summary      Perbarui Tugas
// @Description  Memperbarui judul, deskripsi, status, atau prioritas tugas
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      string                    true  "UUID Tugas"
// @Param        request body      service.UpdateTaskRequest true  "Data Pembaruan"
// @Success      200     {object}  utils.DefaultResponse
// @Failure      400     {object}  utils.DefaultResponse
// @Failure      401     {object}  utils.DefaultResponse
// @Router       /tasks/{id} [put]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Autentikasi gagal", "User tidak valid"))
		return
	}

	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Format ID tugas tidak valid", err.Error()))
		return
	}

	var req service.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Menggunakan penerjemah validasi
		formattedErrors := utils.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Validasi input gagal", formattedErrors))
		return
	}

	task, err := h.taskService.UpdateTask(taskID, userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Gagal memperbarui tugas", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Tugas berhasil diperbarui", task))
}

// DeleteTask godoc
// @Summary      Hapus Tugas
// @Description  Menghapus tugas secara permanen (soft delete)
// @Tags         Tasks
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "UUID Tugas"
// @Success      200  {object}  utils.DefaultResponse
// @Failure      400  {object}  utils.DefaultResponse
// @Failure      401  {object}  utils.DefaultResponse
// @Router       /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Autentikasi gagal", "User tidak valid"))
		return
	}

	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Format ID tugas tidak valid", err.Error()))
		return
	}

	err = h.taskService.DeleteTask(taskID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Gagal menghapus tugas", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Tugas berhasil dihapus", nil))
}

// AddSubTasks godoc
// @Summary      Tambahkan Sub-Task
// @Description  Menambahkan daftar sub-task ke dalam tugas utama (sangat cocok untuk menyimpan hasil AI)
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      string                    true  "UUID Tugas Utama"
// @Param        request body      []service.CreateSubTaskRequest true  "Array Sub Tugas"
// @Success      201     {object}  utils.DefaultResponse
// @Failure      400     {object}  utils.DefaultResponse
// @Failure      401     {object}  utils.DefaultResponse
// @Router       /tasks/{id}/subtasks [post]
func (h *TaskHandler) AddSubTasks(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Autentikasi gagal", "User tidak valid"))
		return
	}

	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Format ID tugas tidak valid", err.Error()))
		return
	}

	// Menerima array dari struct
	var reqs []service.CreateSubTaskRequest
	if err := c.ShouldBindJSON(&reqs); err != nil {
		formattedErrors := utils.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Validasi input gagal", formattedErrors))
		return
	}

	subTasks, err := h.taskService.AddSubTasks(taskID, userID, reqs)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Gagal menambahkan sub-tugas", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("Sub-tugas berhasil ditambahkan", subTasks))
}

// UploadAttachment godoc
// @Summary      Unggah Lampiran Tugas
// @Description  Mengunggah file dokumen (PDF) atau gambar (JPG, PNG) maksimal 5MB ke dalam tugas spesifik
// @Tags         Tasks
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        id        path      string  true  "UUID Tugas"
// @Param        attachment formData  file    true  "File Lampiran"
// @Success      200       {object}  utils.DefaultResponse
// @Failure      400       {object}  utils.DefaultResponse
// @Failure      401       {object}  utils.DefaultResponse
// @Router       /tasks/{id}/upload [post]
func (h *TaskHandler) UploadAttachment(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Autentikasi gagal", "User tidak valid"))
		return
	}

	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Format ID tugas tidak valid", err.Error()))
		return
	}

	// 1. Terima File dari form-data dengan key "attachment"
	file, err := c.FormFile("attachment")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Gagal membaca file", "Pastikan anda mengirim file dengan key 'attachment'"))
		return
	}

	// 2. Proteksi 1: Validasi Ukuran (Maksimal 5 MB)
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Validasi file gagal", "Ukuran file tidak boleh melebihi 5MB"))
		return
	}

	// 3. Proteksi 2: Validasi Ekstensi / Tipe File
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".pdf": true,
	}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Validasi file gagal", "Hanya file JPG, PNG, dan PDF yang diizinkan"))
		return
	}

	// 4. Proteksi 3: Sanitasi Nama File (Menggunakan UUID agar tidak bentrok/ter-overwrite)
	newFileName := uuid.New().String() + ext
	uploadDir := "uploads/attachments"

	// Memastikan folder uploads/attachments tersedia di server fisik
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Gagal menyiapkan direktori penyimpanan", err.Error()))
		return
	}

	// Menentukan path lokasi penyimpanan file fisik
	savePath := filepath.Join(uploadDir, newFileName)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Gagal menyimpan file ke server", err.Error()))
		return
	}

	// 5. Menyimpan URL File ke Database (Bukan file aslinya)
	// Kita simpan format URL statis yang bisa diakses via browser nantinya
	fileURL := "/uploads/attachments/" + newFileName
	task, err := h.taskService.SaveAttachment(taskID, userID, fileURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Gagal memperbarui data tugas", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Lampiran berhasil diunggah", task))
}

// AddCollaborator godoc
// @Summary      Tambah Kolaborator
// @Description  Menambahkan pengguna lain ke dalam tugas dengan akses read_only atau edit
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      string                       true  "UUID Tugas Utama"
// @Param        request body      service.AddCollaboratorRequest true  "Data Kolaborator"
// @Success      201     {object}  utils.DefaultResponse
// @Failure      400     {object}  utils.DefaultResponse
// @Failure      401     {object}  utils.DefaultResponse
// @Router       /tasks/{id}/collaborators [post]
func (h *TaskHandler) AddCollaborator(c *gin.Context) {
	ownerID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Autentikasi gagal", "User tidak valid"))
		return
	}

	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Format ID tugas tidak valid", err.Error()))
		return
	}

	var req service.AddCollaboratorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		formattedErrors := utils.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Validasi input gagal", formattedErrors))
		return
	}

	if err := h.taskService.AddCollaborator(taskID, ownerID, req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Gagal menambahkan kolaborator", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("Kolaborator berhasil ditambahkan", nil))
}