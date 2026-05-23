package repository

import (
	"time"

	"go-reminder/internal/models"
	"gorm.io/gorm"
)

type TaskRepository interface {
	List() ([]models.Task, error)
	FindDueTasks(status models.TaskStatus, dueFrom time.Time, dueTo time.Time) ([]models.Task, error)
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) List() ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.Order("due_at asc").Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) FindDueTasks(status models.TaskStatus, dueFrom time.Time, dueTo time.Time) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.
		Where("status = ? AND due_at BETWEEN ? AND ?", status, dueFrom, dueTo).
		Order("due_at asc").
		Find(&tasks).Error
	return tasks, err
}
