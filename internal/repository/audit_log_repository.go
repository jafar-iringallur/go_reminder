package repository

import (
	"time"

	"go-reminder/internal/models"
	"gorm.io/gorm"
)

type AuditLogRepository interface {
	Create(log *models.AuditLog) error
	List(eventType string, limit int) ([]models.AuditLog, error)
	HasReminderSince(ruleID uint, taskID uint, since time.Time) (bool, error)
}

type auditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Create(log *models.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *auditLogRepository) List(eventType string, limit int) ([]models.AuditLog, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	query := r.db.Order("created_at desc").Limit(limit)
	if eventType != "" {
		query = query.Where("event_type = ?", eventType)
	}

	var logs []models.AuditLog
	err := query.Find(&logs).Error
	return logs, err
}

func (r *auditLogRepository) HasReminderSince(ruleID uint, taskID uint, since time.Time) (bool, error) {
	var count int64
	err := r.db.Model(&models.AuditLog{}).
		Where("event_type = ? AND rule_id = ? AND task_id = ? AND created_at >= ?", models.AuditEventReminderTriggered, ruleID, taskID, since).
		Count(&count).Error
	return count > 0, err
}
