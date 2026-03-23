package repository

import (
	"errors"
	"task-manager-api/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TaskRepository mendefinisikan kontrak fungsi untuk manajemen tugas
type TaskRepository interface {
	Create(task *models.Task) error
	// Tambahan parameter: page, limit, status, priority
	FindAllByUserID(userID uuid.UUID, page, limit int, status, priority string) ([]models.Task, int64, error)
	FindByIDAndUserID(taskID, userID uuid.UUID) (*models.Task, error)
	Update(task *models.Task) error
	Delete(task *models.Task) error
	CreateSubTasks(subTasks []models.SubTask) error
	UpdateAttachment(taskID, userID uuid.UUID, fileURL string) (*models.Task, error)
	UpdateOverdueTasks() (int64, error)
	FindUserByEmail(email string) (*models.User, error)
	AddCollaborator(taskID, userID uuid.UUID, accessLevel string) error
	GetUserRoleInTask(taskID, userID uuid.UUID) (string, error)
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db}
}

func (r *taskRepository) Create(task *models.Task) error {
	return r.db.Create(task).Error
}

// 1. Perbarui fungsi Get All Tasks
func (r *taskRepository) FindAllByUserID(userID uuid.UUID, page, limit int, status, priority string) ([]models.Task, int64, error) {
	var tasks []models.Task
	var totalRows int64

	// Modifikasi Query: Cari tugas di mana user adalah OWNER (user_id = ?) 
	// ATAU user adalah KOLABORATOR (berada di dalam tabel task_collaborators)
	query := r.db.Model(&models.Task{}).
		Where("user_id = ? OR id IN (SELECT task_id FROM task_collaborators WHERE user_id = ?)", userID, userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if priority != "" {
		query = query.Where("priority = ?", priority)
	}

	query.Count(&totalRows)

	offset := (page - 1) * limit
	
	// Tambahkan .Preload("Collaborators") agar array user ikut terbawa di JSON
	err := query.Preload("SubTasks").Preload("Collaborators").
		Offset(offset).Limit(limit).Order("created_at desc").Find(&tasks).Error

	return tasks, totalRows, err
}

// 2. Perbarui fungsi Get Task By ID
func (r *taskRepository) FindByIDAndUserID(taskID, userID uuid.UUID) (*models.Task, error) {
	var task models.Task
	
	// Lakukan modifikasi logika yang sama persis untuk pencarian data tunggal
	err := r.db.Preload("SubTasks").Preload("Collaborators").
		Where("id = ? AND (user_id = ? OR id IN (SELECT task_id FROM task_collaborators WHERE user_id = ?))", taskID, userID, userID).
		First(&task).Error
		
	return &task, err
}

func (r *taskRepository) Update(task *models.Task) error {
	// GORM otomatis menyimpan perubahan berdasarkan ID yang ada di struct task
	return r.db.Save(task).Error
}

func (r *taskRepository) Delete(task *models.Task) error {
	// Karena kita memakai gorm.DeletedAt di model, ini otomatis menjadi Soft Delete
	return r.db.Delete(task).Error
}

func (r *taskRepository) CreateSubTasks(subTasks []models.SubTask) error {
	return r.db.Create(&subTasks).Error // GORM otomatis melakukan Bulk Insert jika diberikan array/slice
}

func (r *taskRepository) UpdateAttachment(taskID, userID uuid.UUID, fileURL string) (*models.Task, error) {
	var task models.Task
	
	// Cari tugasnya dulu untuk memastikan kepemilikan (Data Isolation)
	if err := r.db.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		return nil, err
	}

	// Perbarui kolom attachment dan simpan
	task.AttachmentURL = fileURL
	if err := r.db.Save(&task).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *taskRepository) UpdateOverdueTasks() (int64, error) {
	// Ambil waktu saat ini persis saat fungsi dieksekusi
	now := time.Now()

	// Query: UPDATE tasks SET status = 'overdue' 
	// WHERE due_date < NOW() AND status IN ('pending', 'in_progress') AND deleted_at IS NULL
	result := r.db.Model(&models.Task{}).
		Where("due_date < ? AND status IN ?", now, []string{"pending", "in_progress"}).
		Update("status", "overdue")

	// Mengembalikan jumlah baris data yang berhasil diubah
	return result.RowsAffected, result.Error
}

// Mencari User ID berdasarkan email yang diinputkan
func (r *taskRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Mengeksekusi raw SQL ke tabel pivot. Menggunakan mode UPSERT (ON CONFLICT DO UPDATE)
// agar jika user sudah ada, level aksesnya bisa langsung di-update tanpa error duplikat.
func (r *taskRepository) AddCollaborator(taskID, userID uuid.UUID, accessLevel string) error {
	query := `
		INSERT INTO task_collaborators (task_id, user_id, access_level) 
		VALUES (?, ?, ?) 
		ON CONFLICT (task_id, user_id) 
		DO UPDATE SET access_level = EXCLUDED.access_level, added_at = CURRENT_TIMESTAMP
	`
	return r.db.Exec(query, taskID, userID, accessLevel).Error
}

// Menyelidiki apakah user adalah "owner", "edit", "read_only", atau tidak punya akses
func (r *taskRepository) GetUserRoleInTask(taskID, userID uuid.UUID) (string, error) {
	var task models.Task
	
	// 1. Cek apakah user adalah Pemilik Utama (Owner)
	if err := r.db.Select("user_id").Where("id = ?", taskID).First(&task).Error; err == nil {
		if task.UserID == userID {
			return "owner", nil
		}
	}

	// 2. Jika bukan Owner, cek di tabel pivot (Collaborator)
	var accessLevel string
	err := r.db.Table("task_collaborators").
		Select("access_level").
		Where("task_id = ? AND user_id = ?", taskID, userID).
		Scan(&accessLevel).Error

	if err != nil || accessLevel == "" {
		return "", errors.New("anda tidak memiliki akses ke tugas ini")
	}

	return accessLevel, nil
}