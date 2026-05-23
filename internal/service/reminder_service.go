package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"go-reminder/internal/models"
	"go-reminder/internal/repository"
	"gorm.io/gorm"
)

type ReminderRuleInput struct {
	Name                  string            `json:"name" binding:"required"`
	TaskStatus            models.TaskStatus `json:"task_status" binding:"required"`
	RemindBeforeDueMinute int               `json:"remind_before_due_minutes" binding:"required,min=1"`
	RepeatIntervalMinute  int               `json:"repeat_interval_minutes" binding:"omitempty,min=1"`
	MessageTemplate       string            `json:"message_template" binding:"required"`
	IsActive              *bool             `json:"is_active"`
}

type ReminderService interface {
	CreateRule(input ReminderRuleInput) (*models.ReminderRule, error)
	UpdateRule(id uint, input ReminderRuleInput) (*models.ReminderRule, error)
	DeleteRule(id uint) error
	ListRules() ([]models.ReminderRule, error)
	SetRuleStatus(id uint, active bool) (*models.ReminderRule, error)
	ProcessDueReminders() error
}

type reminderService struct {
	ruleRepo  repository.ReminderRuleRepository
	taskRepo  repository.TaskRepository
	auditRepo repository.AuditLogRepository
}

func NewReminderService(
	ruleRepo repository.ReminderRuleRepository,
	taskRepo repository.TaskRepository,
	auditRepo repository.AuditLogRepository,
) ReminderService {
	return &reminderService{ruleRepo: ruleRepo, taskRepo: taskRepo, auditRepo: auditRepo}
}

func (s *reminderService) CreateRule(input ReminderRuleInput) (*models.ReminderRule, error) {
	if err := validateRuleInput(input); err != nil {
		return nil, err
	}

	active := true
	if input.IsActive != nil {
		active = *input.IsActive
	}

	repeatInterval := input.RepeatIntervalMinute
	if repeatInterval == 0 {
		repeatInterval = 60
	}

	rule := &models.ReminderRule{
		Name:                  strings.TrimSpace(input.Name),
		TaskStatus:            input.TaskStatus,
		RemindBeforeDueMinute: input.RemindBeforeDueMinute,
		RepeatIntervalMinute:  repeatInterval,
		MessageTemplate:       strings.TrimSpace(input.MessageTemplate),
		IsActive:              active,
	}

	if err := s.ruleRepo.Create(rule); err != nil {
		return nil, err
	}

	_ = s.writeAudit(models.AuditEventRuleCreated, "reminder_rule", rule.ID, &rule.ID, nil, rule)
	return rule, nil
}

func (s *reminderService) UpdateRule(id uint, input ReminderRuleInput) (*models.ReminderRule, error) {
	if err := validateRuleInput(input); err != nil {
		return nil, err
	}

	rule, err := s.ruleRepo.FindByID(id)
	if err != nil {
		return nil, normalizeNotFound(err)
	}

	rule.Name = strings.TrimSpace(input.Name)
	rule.TaskStatus = input.TaskStatus
	rule.RemindBeforeDueMinute = input.RemindBeforeDueMinute
	rule.RepeatIntervalMinute = input.RepeatIntervalMinute
	if rule.RepeatIntervalMinute == 0 {
		rule.RepeatIntervalMinute = 60
	}
	rule.MessageTemplate = strings.TrimSpace(input.MessageTemplate)
	if input.IsActive != nil {
		rule.IsActive = *input.IsActive
	}

	if err := s.ruleRepo.Update(rule); err != nil {
		return nil, err
	}

	_ = s.writeAudit(models.AuditEventRuleUpdated, "reminder_rule", rule.ID, &rule.ID, nil, rule)
	return rule, nil
}

func (s *reminderService) DeleteRule(id uint) error {
	rule, err := s.ruleRepo.FindByID(id)
	if err != nil {
		return normalizeNotFound(err)
	}

	if err := s.ruleRepo.Delete(rule); err != nil {
		return err
	}

	return s.writeAudit(models.AuditEventRuleDeleted, "reminder_rule", id, &id, nil, rule)
}

func (s *reminderService) ListRules() ([]models.ReminderRule, error) {
	return s.ruleRepo.List()
}

func (s *reminderService) SetRuleStatus(id uint, active bool) (*models.ReminderRule, error) {
	rule, err := s.ruleRepo.FindByID(id)
	if err != nil {
		return nil, normalizeNotFound(err)
	}

	rule.IsActive = active
	if err := s.ruleRepo.Update(rule); err != nil {
		return nil, err
	}

	_ = s.writeAudit(models.AuditEventRuleStatusChanged, "reminder_rule", id, &id, nil, map[string]any{
		"rule_id":   id,
		"is_active": active,
	})
	return rule, nil
}

func (s *reminderService) ProcessDueReminders() error {
	rules, err := s.ruleRepo.ListActive()
	if err != nil {
		return err
	}

	now := time.Now()
	for _, rule := range rules {
		dueTo := now.Add(time.Duration(rule.RemindBeforeDueMinute) * time.Minute)
		tasks, err := s.taskRepo.FindDueTasks(rule.TaskStatus, now, dueTo)
		if err != nil {
			return err
		}

		for _, task := range tasks {
			alreadySent, err := s.auditRepo.HasReminderSince(rule.ID, task.ID, now.Add(-time.Duration(rule.RepeatIntervalMinute)*time.Minute))
			if err != nil {
				return err
			}
			if alreadySent {
				continue
			}

			message := renderReminderMessage(rule.MessageTemplate, rule, task)
			log.Printf("[REMINDER] rule=%d task=%d message=%q due_at=%s", rule.ID, task.ID, message, task.DueAt.Format(time.RFC3339))

			if err := s.writeAudit(models.AuditEventReminderTriggered, "task", task.ID, &rule.ID, &task.ID, map[string]any{
				"rule":    rule,
				"task":    task,
				"message": message,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateRuleInput(input ReminderRuleInput) error {
	if strings.TrimSpace(input.Name) == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidInput)
	}
	if strings.TrimSpace(input.MessageTemplate) == "" {
		return fmt.Errorf("%w: message_template is required", ErrInvalidInput)
	}
	if input.TaskStatus != models.TaskStatusPending &&
		input.TaskStatus != models.TaskStatusInProgress &&
		input.TaskStatus != models.TaskStatusCompleted {
		return fmt.Errorf("%w: task_status must be one of: pending, in_progress, completed", ErrInvalidInput)
	}
	if input.RemindBeforeDueMinute <= 0 {
		return fmt.Errorf("%w: remind_before_due_minutes must be greater than zero", ErrInvalidInput)
	}
	if input.RepeatIntervalMinute < 0 {
		return fmt.Errorf("%w: repeat_interval_minutes cannot be negative", ErrInvalidInput)
	}
	return nil
}

func (s *reminderService) writeAudit(eventType models.AuditEventType, entity string, entityID uint, ruleID *uint, taskID *uint, details any) error {
	payload, err := json.Marshal(details)
	if err != nil {
		return err
	}

	return s.auditRepo.Create(&models.AuditLog{
		EventType: eventType,
		Entity:    entity,
		EntityID:  entityID,
		RuleID:    ruleID,
		TaskID:    taskID,
		Details:   models.JSONDetails(payload),
	})
}

func normalizeNotFound(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	return err
}

func renderReminderMessage(template string, rule models.ReminderRule, task models.Task) string {
	replacements := map[string]string{
		"{{rule_name}}":   rule.Name,
		"{{task_title}}":  task.Title,
		"{{task_status}}": string(task.Status),
		"{{due_at}}":      task.DueAt.Format(time.RFC3339),
	}

	message := template
	for token, value := range replacements {
		message = strings.ReplaceAll(message, token, value)
	}

	if message == template {
		return fmt.Sprintf("%s: %s is due at %s", rule.Name, task.Title, task.DueAt.Format(time.RFC3339))
	}
	return message
}
