package service

import (
	"task-manager-api/internal/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// =========================================================================
// 1. SETUP MOCK REPOSITORY FULL UPGRADE
// =========================================================================
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(task *models.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepository) FindAllByUserID(userID uuid.UUID, page, limit int, status, priority string) ([]models.Task, int64, error) {
	args := m.Called(userID, page, limit, status, priority)
	return args.Get(0).([]models.Task), args.Get(1).(int64), args.Error(2)
}

func (m *MockTaskRepository) FindByIDAndUserID(taskID, userID uuid.UUID) (*models.Task, error) {
	args := m.Called(taskID, userID)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskRepository) Update(task *models.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepository) Delete(task *models.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

// Tambahan Mock untuk fitur baru
func (m *MockTaskRepository) CreateSubTasks(subTasks []models.SubTask) error {
	args := m.Called(subTasks)
	return args.Error(0)
}

func (m *MockTaskRepository) UpdateAttachment(taskID, userID uuid.UUID, fileURL string) (*models.Task, error) {
	args := m.Called(taskID, userID, fileURL)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskRepository) FindUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskRepository) AddCollaborator(taskID, userID uuid.UUID, accessLevel string) error {
	args := m.Called(taskID, userID, accessLevel)
	return args.Error(0)
}

func (m *MockTaskRepository) GetUserRoleInTask(taskID, userID uuid.UUID) (string, error) {
	args := m.Called(taskID, userID)
	return args.String(0), args.Error(1)
}

func (m *MockTaskRepository) UpdateOverdueTasks() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

// =========================================================================
// 2. SKENARIO TEST CREATE TASK
// =========================================================================

func TestCreateTask_Success(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	taskService := NewTaskService(mockRepo)

	userID := uuid.New()
	req := CreateTaskRequest{
		Title:       "Belajar Unit Test Golang",
		Description: "Paham penggunaan Testify",
		Priority:    "high",
	}

	mockRepo.On("Create", mock.AnythingOfType("*models.Task")).Return(nil)

	resultTask, err := taskService.CreateTask(userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, resultTask)
	assert.Equal(t, req.Title, resultTask.Title)
	assert.Equal(t, userID, resultTask.UserID)

	mockRepo.AssertExpectations(t)
}

// =========================================================================
// 3. SKENARIO TEST UPDATE TASK (DENGAN RBAC)
// =========================================================================

func TestUpdateTask_Success(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	taskService := NewTaskService(mockRepo)

	taskID := uuid.New()
	userID := uuid.New()
	req := UpdateTaskRequest{Title: "Judul Baru"}

	existingTask := &models.Task{ID: taskID, UserID: userID, Title: "Judul Lama"}

	// Mocking otorisasi dan pengambilan data
	mockRepo.On("GetUserRoleInTask", taskID, userID).Return("owner", nil)
	mockRepo.On("FindByIDAndUserID", taskID, userID).Return(existingTask, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Task")).Return(nil)

	result, err := taskService.UpdateTask(taskID, userID, req)

	assert.NoError(t, err)
	assert.Equal(t, "Judul Baru", result.Title)
	mockRepo.AssertExpectations(t)
}

func TestUpdateTask_ReadOnlyBlocked(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	taskService := NewTaskService(mockRepo)

	taskID := uuid.New()
	userID := uuid.New()
	req := UpdateTaskRequest{Title: "Hekel Beraksi"}

	// Mocking simulasi user sebagai read_only
	mockRepo.On("GetUserRoleInTask", taskID, userID).Return("read_only", nil)

	result, err := taskService.UpdateTask(taskID, userID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "read-only")
	// Kita tidak perlu assert .On("Update") karena fungsi harus terhenti sebelum itu

	mockRepo.AssertExpectations(t)
}

// =========================================================================
// 4. SKENARIO TEST DELETE TASK (EKSKLUSIF OWNER)
// =========================================================================

func TestDeleteTask_EditorBlocked(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	taskService := NewTaskService(mockRepo)

	taskID := uuid.New()
	userID := uuid.New()

	// Simulasi user adalah editor, bukan owner
	mockRepo.On("GetUserRoleInTask", taskID, userID).Return("edit", nil)

	err := taskService.DeleteTask(taskID, userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "hanya pemilik utama")
	mockRepo.AssertExpectations(t)
}

// =========================================================================
// 5. SKENARIO TEST KOLABORATOR (AddCollaborator)
// =========================================================================

func TestAddCollaborator_SelfInviteBlocked(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	taskService := NewTaskService(mockRepo)

	taskID := uuid.New()
	ownerID := uuid.New() // ID yang sama dengan invitee
	req := AddCollaboratorRequest{
		Email:       "me@example.com",
		AccessLevel: "edit",
	}

	invitee := &models.User{ID: ownerID, Email: "me@example.com"}

	// 1. Validasi role owner lewat
	mockRepo.On("GetUserRoleInTask", taskID, ownerID).Return("owner", nil)
	// 2. Pencarian email berhasil
	mockRepo.On("FindUserByEmail", req.Email).Return(invitee, nil)

	err := taskService.AddCollaborator(taskID, ownerID, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "menambahkan diri sendiri")

	mockRepo.AssertExpectations(t)
}