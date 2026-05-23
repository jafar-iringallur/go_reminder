package repository

import (
	"go-reminder/internal/models"
	"gorm.io/gorm"
)

type ReminderRuleRepository interface {
	Create(rule *models.ReminderRule) error
	Update(rule *models.ReminderRule) error
	Delete(rule *models.ReminderRule) error
	FindByID(id uint) (*models.ReminderRule, error)
	List() ([]models.ReminderRule, error)
	ListActive() ([]models.ReminderRule, error)
}

type reminderRuleRepository struct {
	db *gorm.DB
}

func NewReminderRuleRepository(db *gorm.DB) ReminderRuleRepository {
	return &reminderRuleRepository{db: db}
}

func (r *reminderRuleRepository) Create(rule *models.ReminderRule) error {
	return r.db.Create(rule).Error
}

func (r *reminderRuleRepository) Update(rule *models.ReminderRule) error {
	return r.db.Save(rule).Error
}

func (r *reminderRuleRepository) Delete(rule *models.ReminderRule) error {
	return r.db.Delete(rule).Error
}

func (r *reminderRuleRepository) FindByID(id uint) (*models.ReminderRule, error) {
	var rule models.ReminderRule
	if err := r.db.First(&rule, id).Error; err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *reminderRuleRepository) List() ([]models.ReminderRule, error) {
	var rules []models.ReminderRule
	err := r.db.Order("id desc").Find(&rules).Error
	return rules, err
}

func (r *reminderRuleRepository) ListActive() ([]models.ReminderRule, error) {
	var rules []models.ReminderRule
	err := r.db.Where("is_active = ?", true).Order("id asc").Find(&rules).Error
	return rules, err
}
