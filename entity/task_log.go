package entity

import "time"

type TaskLog struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	TaskID         uint      `gorm:"not null" json:"task_id"`
	PerformedBy    uint      `gorm:"not null" json:"performed_by"`
	PreviousUserID uint      `gorm:"not null" json:"previous_user_id"`
	NewUserID      uint      `gorm:"not null" json:"new_user_id"`
	Action         string    `gorm:"type:varchar(50);not null" json:"action"`
	CreatedAt      time.Time `json:"created_at"`
}
