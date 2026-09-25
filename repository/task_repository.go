package repository

import (
	"flow-desk/dto"
	"flow-desk/entity"

	"gorm.io/gorm"
)

type TaskRepository interface {
	Create(task *entity.Task) error
	FindByUserID(userID uint, query dto.TaskQueryParam) ([]entity.Task, int64, error)
	FindByIDAndUserID(id uint, userID uint) (*entity.Task, error)
	Update(task *entity.Task) error
	Delete(id uint, userID uint) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(task *entity.Task) error {
	return r.db.Create(task).Error
}

func (r *taskRepository) FindByUserID(userID uint, query dto.TaskQueryParam) ([]entity.Task, int64, error) {
	var tasks []entity.Task
	var totalData int64
	// Base query (wajib filter user_id agar terisolasi)
	dbQuery := r.db.Model(&entity.Task{}).Where("user_id = ?", userID)
	// Filter berdasarkan Status (jika ada)
	if query.Status != "" {
		dbQuery = dbQuery.Where("status = ?", query.Status)
	}
	// Filter berdasarkan Search pada Title atau Description (jika ada)
	if query.Search != "" {
		searchTerm := "%" + query.Search + "%"
		dbQuery = dbQuery.Where("title LIKE ? OR description LIKE ?", searchTerm, searchTerm)
	}
	// Hitung total data sebelum dilakukan Limit & Offset
	err := dbQuery.Count(&totalData).Error
	if err != nil {
		return nil, 0, err
	}
	// Hitung Offset untuk Pagination
	offset := (query.Page - 1) * query.Limit
	// Eksekusi query dengan Order, Limit, dan Offset
	err = dbQuery.Order("created_at desc").Limit(query.Limit).Offset(offset).Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}
	return tasks, totalData, nil
}

func (r *taskRepository) FindByIDAndUserID(id uint, userID uint) (*entity.Task, error) {
	var task entity.Task
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) Update(task *entity.Task) error {
	return r.db.Save(task).Error
}

func (r *taskRepository) Delete(id uint, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&entity.Task{}).Error
}
