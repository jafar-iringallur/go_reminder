package service

import (
	"go-reminder/internal/models"
	"go-reminder/internal/repository"
)

type AuditService interface {
	List(eventType string, limit int) ([]models.AuditLog, error)
}

type auditService struct {
	auditRepo repository.AuditLogRepository
}

func NewAuditService(auditRepo repository.AuditLogRepository) AuditService {
	return &auditService{auditRepo: auditRepo}
}

func (s *auditService) List(eventType string, limit int) ([]models.AuditLog, error) {
	return s.auditRepo.List(eventType, limit)
}
