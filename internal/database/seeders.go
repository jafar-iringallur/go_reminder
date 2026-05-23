package database

import (
	"time"

	"go-reminder/internal/models"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	if err := SeedTasks(db); err != nil {
		return err
	}
	return SeedReminderRules(db)
}

func SeedTasks(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.Task{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	now := time.Now()
	tasks := []models.Task{
		{Title: "Submit API machine test", Description: "Finalize Gin reminder system implementation.", Status: models.TaskStatusPending, DueAt: now.Add(10 * time.Minute)},
		{Title: "Review database migrations", Description: "Confirm schema and seed data are working.", Status: models.TaskStatusInProgress, DueAt: now.Add(25 * time.Minute)},
		{Title: "Prepare deployment notes", Description: "Document run steps and environment variables.", Status: models.TaskStatusPending, DueAt: now.Add(55 * time.Minute)},
		{Title: "Archive completed onboarding checklist", Description: "Already completed sample task.", Status: models.TaskStatusCompleted, DueAt: now.Add(90 * time.Minute)},
		{Title: "Follow up on QA feedback", Description: "Pending task due later today.", Status: models.TaskStatusPending, DueAt: now.Add(3 * time.Hour)},
	}

	return db.Create(&tasks).Error
}

func SeedReminderRules(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.ReminderRule{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	rule := models.ReminderRule{
		Name:                  "Pending tasks due within one hour",
		TaskStatus:            models.TaskStatusPending,
		RemindBeforeDueMinute: 60,
		RepeatIntervalMinute:  15,
		MessageTemplate:       "Reminder: {{task_title}} is due at {{due_at}}",
		IsActive:              true,
	}

	return db.Create(&rule).Error
}
