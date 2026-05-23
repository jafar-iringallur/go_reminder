package models

import "time"

type ReminderRule struct {
	ID                    uint       `json:"id" gorm:"primaryKey"`
	Name                  string     `json:"name" gorm:"not null"`
	TaskStatus            TaskStatus `json:"task_status" gorm:"type:varchar(30);index;not null"`
	RemindBeforeDueMinute int        `json:"remind_before_due_minutes" gorm:"not null"`
	RepeatIntervalMinute  int        `json:"repeat_interval_minutes" gorm:"not null;default:60"`
	MessageTemplate       string     `json:"message_template" gorm:"not null"`
	IsActive              bool       `json:"is_active" gorm:"index;not null;default:true"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}
