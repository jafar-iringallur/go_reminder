package service

import (
	"go-reminder/internal/models"
	"go-reminder/internal/repository"
)

type TaskService interface {
	List() ([]models.Task, error)
}

type taskService struct {
	taskRepo repository.TaskRepository
}

func NewTaskService(taskRepo repository.TaskRepository) TaskService {
	return &taskService{taskRepo: taskRepo}
}

func (s *taskService) List() ([]models.Task, error) {
	return s.taskRepo.List()
}
