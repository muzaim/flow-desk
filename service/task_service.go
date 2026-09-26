package service

import (
	"math"

	"flow-desk/dto"
	"flow-desk/entity"
	"flow-desk/repository"
	"flow-desk/utils"
)

type TaskService interface {
	CreateTask(userID uint, req dto.CreateTaskRequest) (*dto.TaskResponse, error)
	GetTasks(userID uint, query dto.TaskQueryParam) (*dto.PaginatedTaskResponse, error)
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
		return nil, utils.NewInternalServerError("CREATE_TASK_ERROR", "Gagal membuat task", err)
	}

	return formatTaskResponse(&task), nil
}

func (s *taskService) GetTasks(userID uint, query dto.TaskQueryParam) (*dto.PaginatedTaskResponse, error) {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Limit <= 0 {
		query.Limit = 10
	}
	tasks, totalData, err := s.taskRepo.FindByUserID(userID, query)
	if err != nil {
		return nil, utils.NewInternalServerError("GET_TASKS_ERROR", "Gagal mengambil data task", err)
	}
	var taskResponses []dto.TaskResponse
	for _, task := range tasks {
		taskResponses = append(taskResponses, *formatTaskResponse(&task))
	}
	totalPage := int(math.Ceil(float64(totalData) / float64(query.Limit)))
	return &dto.PaginatedTaskResponse{
		Data: taskResponses,
		Meta: dto.MetaResponse{
			TotalData:   totalData,
			TotalPage:   totalPage,
			CurrentPage: query.Page,
			Limit:       query.Limit,
		},
	}, nil
}

func (s *taskService) GetTaskByID(id uint, userID uint) (*dto.TaskResponse, error) {
	task, err := s.taskRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return nil, utils.NewNotFoundError("TASK_NOT_FOUND", "Task tidak ditemukan", err)
	}

	return formatTaskResponse(task), nil
}

func (s *taskService) UpdateTask(id uint, userID uint, req dto.UpdateTaskRequest) (*dto.TaskResponse, error) {
	task, err := s.taskRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return nil, utils.NewNotFoundError("TASK_NOT_FOUND", "Task tidak ditemukan", err)
	}

	task.Title = req.Title
	task.Description = req.Description
	task.Status = req.Status

	err = s.taskRepo.Update(task)
	if err != nil {
		return nil, utils.NewInternalServerError("UPDATE_TASK_ERROR", "Gagal memperbarui task", err)
	}

	return formatTaskResponse(task), nil
}

func (s *taskService) DeleteTask(id uint, userID uint) error {
	_, err := s.taskRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return utils.NewNotFoundError("TASK_NOT_FOUND", "Task tidak ditemukan", err)
	}

	err = s.taskRepo.Delete(id, userID)
	if err != nil {
		return utils.NewInternalServerError("DELETE_TASK_ERROR", "Gagal menghapus task", err)
	}

	return nil
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
