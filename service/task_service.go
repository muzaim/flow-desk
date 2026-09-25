package service

import (
	"errors"
	"math"

	"flow-desk/dto"
	"flow-desk/entity"
	"flow-desk/repository"
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

	var tags []entity.TaskTag
	for _, tagName := range req.Tags {
		tags = append(tags, entity.TaskTag{Name: tagName})
	}

	task := entity.Task{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		Tags:        tags,
	}

	err := s.taskRepo.Create(&task)
	if err != nil {
		return nil, errors.New("gagal membuat task")
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
		return nil, errors.New("gagal mengambil data task")
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
		return nil, errors.New("task tidak ditemukan")
	}

	return formatTaskResponse(task), nil
}

func (s *taskService) UpdateTask(id uint, userID uint, req dto.UpdateTaskRequest) (*dto.TaskResponse, error) {
	task, err := s.taskRepo.FindByIDAndUserID(id, userID)

	if err != nil {
		return nil, errors.New("task tidak ditemukan")
	}

	var newTags []entity.TaskTag
	for _, tagName := range req.Tags {
		newTags = append(newTags, entity.TaskTag{TaskID: id, Name: tagName})
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
	var tagNames []string
	for _, tag := range task.Tags {
		tagNames = append(tagNames, tag.Name)
	}
	if tagNames == nil {
		tagNames = []string{}
	}
	return &dto.TaskResponse{
		ID:          task.ID,
		UserID:      task.UserID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Tags:        tagNames,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}
