package service

import (
	"errors"
	"math"
	"task-manager-api/internal/models"
	"task-manager-api/internal/repository"
	"time"

	"github.com/google/uuid"
)

// DTO untuk pembuatan tugas baru
type CreateTaskRequest struct {
	Title       string     `json:"title" binding:"required,max=150"`
	Description string     `json:"description"`
	Priority    string     `json:"priority" binding:"omitempty,oneof=low medium high"`
	DueDate     *time.Time `json:"due_date"`
}

type CreateSubTaskRequest struct {
	Title       string `json:"title" binding:"required,max=150"`
	Description string `json:"description"`
	Priority    string `json:"priority" binding:"omitempty,oneof=low medium high"`
}

// DTO untuk pembaruan tugas (Perbaikan: Menambahkan 'overdue' untuk Cron Job)
type UpdateTaskRequest struct {
	Title       string     `json:"title" binding:"omitempty,max=150"`
	Description string     `json:"description"`
	Status      string     `json:"status" binding:"omitempty,oneof=pending in_progress completed overdue"`
	Priority    string     `json:"priority" binding:"omitempty,oneof=low medium high"`
	DueDate     *time.Time `json:"due_date"`
}

// DTO untuk parameter pencarian URL
type TaskFilterParams struct {
	Page     int
	Limit    int
	Status   string
	Priority string
}

// DTO untuk membungkus respons Pagination
type PaginatedResponse struct {
	Data       []models.Task `json:"data"`
	TotalItems int64         `json:"total_items"`
	TotalPages int           `json:"total_pages"`
	Page       int           `json:"current_page"`
	Limit      int           `json:"limit"`
}

type AddCollaboratorRequest struct {
	Email       string `json:"email" binding:"required,email"`
	AccessLevel string `json:"access_level" binding:"required,oneof=read_only edit"`
}

type TaskService interface {
	CreateTask(userID uuid.UUID, req CreateTaskRequest) (*models.Task, error)
	GetTasksByUser(userID uuid.UUID, filter TaskFilterParams) (PaginatedResponse, error)
	GetTaskByID(taskID, userID uuid.UUID) (*models.Task, error)
	UpdateTask(taskID, userID uuid.UUID, req UpdateTaskRequest) (*models.Task, error)
	DeleteTask(taskID, userID uuid.UUID) error
	AddSubTasks(taskID, userID uuid.UUID, reqs []CreateSubTaskRequest) ([]models.SubTask, error)
	SaveAttachment(taskID, userID uuid.UUID, fileURL string) (*models.Task, error)
	AddCollaborator(taskID, ownerID uuid.UUID, req AddCollaboratorRequest) error
}

type taskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{repo}
}

func (s *taskService) CreateTask(userID uuid.UUID, req CreateTaskRequest) (*models.Task, error) {
	priority := req.Priority
	if priority == "" {
		priority = "medium"
	}

	task := &models.Task{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    priority,
		Status:      "pending",
		DueDate:     req.DueDate,
	}

	err := s.repo.Create(task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *taskService) GetTasksByUser(userID uuid.UUID, filter TaskFilterParams) (PaginatedResponse, error) {
	tasks, totalRows, err := s.repo.FindAllByUserID(userID, filter.Page, filter.Limit, filter.Status, filter.Priority)
	if err != nil {
		return PaginatedResponse{}, err
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(filter.Limit)))

	response := PaginatedResponse{
		Data:       tasks,
		TotalItems: totalRows,
		TotalPages: totalPages,
		Page:       filter.Page,
		Limit:      filter.Limit,
	}

	return response, nil
}

func (s *taskService) GetTaskByID(taskID, userID uuid.UUID) (*models.Task, error) {
	return s.repo.FindByIDAndUserID(taskID, userID)
}

func (s *taskService) UpdateTask(taskID, userID uuid.UUID, req UpdateTaskRequest) (*models.Task, error) {
	// Validasi Hak Akses Khusus Mutasi Data
	role, err := s.repo.GetUserRoleInTask(taskID, userID)
	if err != nil {
		return nil, errors.New("tugas tidak ditemukan atau anda tidak memiliki akses")
	}
	if role == "read_only" {
		return nil, errors.New("akses ditolak: anda hanya memiliki hak baca (read-only) untuk tugas ini")
	}

	task, err := s.repo.FindByIDAndUserID(taskID, userID)
	if err != nil {
		return nil, errors.New("tugas tidak ditemukan")
	}

	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.Status != "" {
		task.Status = req.Status
	}
	if req.Priority != "" {
		task.Priority = req.Priority
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}

	err = s.repo.Update(task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *taskService) DeleteTask(taskID, userID uuid.UUID) error {
	// 1. Validasi Hak Akses Eksklusif Pemilik (Authorization)
	role, err := s.repo.GetUserRoleInTask(taskID, userID)
	if err != nil {
		return errors.New("tugas tidak ditemukan atau anda tidak memiliki akses")
	}
	
	// Blokir siapa pun yang BUKAN owner (termasuk editor dan read_only)
	if role != "owner" {
		return errors.New("akses ditolak: hanya pemilik utama (owner) yang dapat menghapus tugas ini")
	}

	// 2. Ambil data object tugas secara utuh dari database
	task, err := s.repo.FindByIDAndUserID(taskID, userID)
	if err != nil {
		return errors.New("tugas tidak ditemukan atau sudah terhapus")
	}

	// 3. Eksekusi penghapusan (Soft Delete) dengan mengirimkan object model
	return s.repo.Delete(task)
}

func (s *taskService) AddSubTasks(taskID, userID uuid.UUID, reqs []CreateSubTaskRequest) ([]models.SubTask, error) {
	// Validasi Hak Akses (Blokir Read-Only)
	role, err := s.repo.GetUserRoleInTask(taskID, userID)
	if err != nil {
		return nil, errors.New("tugas utama tidak ditemukan atau anda tidak memiliki akses")
	}
	if role == "read_only" {
		return nil, errors.New("akses ditolak: anda hanya memiliki hak baca (read-only) untuk tugas ini")
	}

	var subTasks []models.SubTask
	for _, req := range reqs {
		priority := req.Priority
		if priority == "" {
			priority = "medium"
		}

		subTasks = append(subTasks, models.SubTask{
			TaskID:      taskID,
			Title:       req.Title,
			Description: req.Description,
			Priority:    priority,
			Status:      "pending",
		})
	}

	err = s.repo.CreateSubTasks(subTasks)
	if err != nil {
		return nil, err
	}

	return subTasks, nil
}

func (s *taskService) SaveAttachment(taskID, userID uuid.UUID, fileURL string) (*models.Task, error) {
	// Validasi Hak Akses (Blokir Read-Only)
	role, err := s.repo.GetUserRoleInTask(taskID, userID)
	if err != nil {
		return nil, errors.New("tugas utama tidak ditemukan atau anda tidak memiliki akses")
	}
	if role == "read_only" {
		return nil, errors.New("akses ditolak: anda hanya memiliki hak baca (read-only) untuk tugas ini")
	}

	return s.repo.UpdateAttachment(taskID, userID, fileURL)
}

func (s *taskService) AddCollaborator(taskID, ownerID uuid.UUID, req AddCollaboratorRequest) error {
	// Validasi Hak Akses Eksklusif Pemilik (Hanya owner yang bisa mengundang orang lain)
	role, err := s.repo.GetUserRoleInTask(taskID, ownerID)
	if err != nil {
		return errors.New("tugas tidak ditemukan atau anda tidak memiliki akses")
	}
	if role != "owner" {
		return errors.New("akses ditolak: hanya pemilik utama tugas ini yang dapat menambahkan kolaborator")
	}

	invitee, err := s.repo.FindUserByEmail(req.Email)
	if err != nil {
		return errors.New("pengguna dengan email tersebut tidak ditemukan di sistem")
	}

	if invitee.ID == ownerID {
		return errors.New("anda tidak bisa menambahkan diri sendiri sebagai kolaborator")
	}

	return s.repo.AddCollaborator(taskID, invitee.ID, req.AccessLevel)
}