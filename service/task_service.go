package service

import (
	"errors"

	"flow-desk/dto"
	"flow-desk/entity"
	"flow-desk/repository"
)

type TaskService interface {
	CreateTask(userID uint, req dto.CreateTaskRequest) (*dto.TaskResponse, error)
	GetTasks(userID uint) ([]dto.TaskResponse, error)
	GetTaskByID(id uint, userID uint) (*dto.TaskResponse, error)
	UpdateTask(id uint, userID uint, req dto.UpdateTaskRequest) (*dto.TaskResponse, error)
	DeleteTask(id uint, userID uint) error
}

type taskService struct {
	taskRepo repository.TaskRepository
}

func NewTaskService(taskRepo repository.TaskRepository) TaskService {
	return &taskService{taskRepo: taskRepo}
}

func (s *taskService) CreateTask(userID uint, req dto.CreateTaskRequest) (*dto.TaskResponse, error) {
	status := req.Status
	if status == "" {
		status = "pending"
	}

	task := entity.Task{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
	}

	err := s.taskRepo.Create(&task)
	if err != nil {
		return nil, errors.New("gagal membuat task")
	}

	return formatTaskResponse(&task), nil
}

func (s *taskService) GetTasks(userID uint) ([]dto.TaskResponse, error) {
	tasks, err := s.taskRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("gagal mengambil data task")
	}

	var res []dto.TaskResponse
	for _, task := range tasks {
		res = append(res, *formatTaskResponse(&task))
	}
	return res, nil
}

func (s *taskService) GetTaskByID(id uint, userID uint) (*dto.TaskResponse, error) {
	task, err := s.taskRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return nil, errors.New("task tidak ditemukan")
	}

	return formatTaskResponse(task), nil
}

func (s *taskService) UpdateTask(id uint, userID uint, req dto.UpdateTaskRequest) (*dto.TaskResponse, error) {
	task, err := s.taskRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return nil, errors.New("task tidak ditemukan")
	}

	task.Title = req.Title
	task.Description = req.Description
	task.Status = req.Status

	err = s.taskRepo.Update(task)
	if err != nil {
		return nil, errors.New("gagal memperbarui task")
	}

	return formatTaskResponse(task), nil
}

func (s *taskService) DeleteTask(id uint, userID uint) error {
	_, err := s.taskRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return errors.New("task tidak ditemukan")
	}

	return s.taskRepo.Delete(id, userID)
}

func formatTaskResponse(task *entity.Task) *dto.TaskResponse {
	return &dto.TaskResponse{
		ID:          task.ID,
		UserID:      task.UserID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}
