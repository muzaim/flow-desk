package repository

import (
	"flow-desk/entity"

	"gorm.io/gorm"
)

type TaskRepository interface {
	Create(task *entity.Task) error
	FindByUserID(userID uint) ([]entity.Task, error)
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

func (r *taskRepository) FindByUserID(userID uint) ([]entity.Task, error) {
	var tasks []entity.Task
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&tasks).Error
	return tasks, err
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
