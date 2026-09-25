package dto

import "time"

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required,min=3,max=255"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title" binding:"required,min=3,max=255"`
	Description string `json:"description"`
	Status      string `json:"status" binding:"required,oneof=pending in_progress completed"`
}

type TaskResponse struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// (Tambahkan struct ini di file dto/task_dto.go)

type TaskQueryParam struct {
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
	Search string `form:"search"`
	Status string `form:"status"`
}

type MetaResponse struct {
	TotalData   int64 `json:"total_data"`
	TotalPage   int   `json:"total_page"`
	CurrentPage int   `json:"current_page"`
	Limit       int   `json:"limit"`
}

type PaginatedTaskResponse struct {
	Data []TaskResponse `json:"data"`
	Meta MetaResponse   `json:"meta"`
}
